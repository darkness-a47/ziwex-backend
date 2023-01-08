package services

import (
	"net/http"
	"strings"
	"ziwex/db"
	"ziwex/dtos"
	"ziwex/minioClient"
	"ziwex/models"
	"ziwex/types"
	"ziwex/types/emptyResponse"
	"ziwex/types/fileResponse"
	"ziwex/types/jsonResponse"
	"ziwex/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
)

func UploadFile(d dtos.UploadFile) types.Response {
	r := &jsonResponse.Response{}
	file, fileOpenErr := d.File.Open()
	if fileOpenErr != nil {
		r.Error(fileOpenErr)
		return r
	}
	defer file.Close()

	objectName, uuidErr := uuid.NewRandom()
	if uuidErr != nil {
		r.Error(uuidErr)
		return r
	}

	contentType, contentTypeErr := utils.GetFileContentType(file)
	if contentTypeErr != nil {
		r.Error(contentTypeErr)
		return r
	}

	minioCtx, minioCancel := utils.GetMinioContext()
	defer minioCancel()
	info, minioErr := minioClient.Conn.PutObject(minioCtx, minioClient.ImageBucket, objectName.String(),
		file, d.File.Size, minio.PutObjectOptions{ContentType: contentType})
	if minioErr != nil {
		r.Error(minioErr)
		return r
	}

	dbCtx, dbCancel := utils.GetDatabaseContext()
	defer dbCancel()

	dbErr := db.Poll.QueryRow(dbCtx, `--sql
		INSERT INTO files (filename, file_id, hash_md5, content_type) VALUES ($1, $2, $3::UUID, $4);
	`, d.File.Filename, objectName.String(), info.ETag, contentType).Scan()

	//TODO: revert file insert
	if dbErr != nil && dbErr != pgx.ErrNoRows {
		r.Error(dbErr)
		return r
	}

	r.Write(http.StatusCreated, jsonResponse.Json{
		"message": "image created",
		"id":      objectName,
	})
	return r
}

func ServeFile(d dtos.ServeFile) types.Response {
	r := &jsonResponse.Response{}
	dbCtx, dbCancel := utils.GetDatabaseContext()
	defer dbCancel()

	dbFile := models.File{}
	dbErr := db.Poll.QueryRow(dbCtx, `--sql
		SELECT filename, content_type FROM files WHERE file_id = $1;
	`, d.FileId).Scan(&dbFile.Filename, &dbFile.ContentType)
	if dbErr != nil {
		if dbErr == pgx.ErrNoRows {
			r.Write(http.StatusNotFound, jsonResponse.Json{
				"message": "file not found",
			})
			return r
		}
		r.Error(dbErr)
		return r
	}

	minioCtx, minioCancel := utils.GetMinioGetContext()
	file, minioErr := minioClient.Conn.GetObject(minioCtx, minioClient.ImageBucket, d.FileId, minio.GetObjectOptions{})
	if minioErr != nil {
		r.Error(minioErr)
		return r
	}

	rf := &fileResponse.Response{}
	rf.Write(http.StatusOK, &fileResponse.File{
		File:           file,
		ContentType:    dbFile.ContentType,
		CancelFunction: minioCancel,
	})

	return rf
}

func GetFiles(d dtos.GetFiles) types.Response {
	r := &jsonResponse.Response{}

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	skip := (d.Page - 1) * d.DataPerPage
	rows, err := db.Poll.Query(ctx, `--sql
		SELECT id, filename, file_id, content_type, hash_md5, COUNT(*) OVER() AS total FROM files
		WHERE (filename LIKE '%' || $1 || '%') ORDER BY id DESC OFFSET $2 LIMIT $3;
	`, *d.Filename, skip, d.DataPerPage)
	if err != nil {
		r.Error(err)
		return r
	}
	var totalRows int
	files := make([]models.File, 0)
	for rows.Next() {
		f := models.File{}
		err := rows.Scan(&f.Id, &f.Filename, &f.FileId, &f.ContentType, &f.HashMd5, &totalRows)
		if err != nil {
			r.Error(err)
			return r
		}
		f.HashMd5 = strings.ReplaceAll(f.HashMd5, "-", "")
		files = append(files, f)
	}

	r.Write(http.StatusOK, jsonResponse.Json{
		"message":    "ok",
		"files":      files,
		"total_rows": totalRows,
	})

	return r
}

func UpdateFile(d dtos.UpdateFile) types.Response {
	r := &jsonResponse.Response{}

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()
	f := models.File{}
	err := db.Poll.QueryRow(ctx, `--sql
		SELECT file_id FROM files WHERE id = $1;
	`, d.Id).Scan(&f.FileId)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.Write(http.StatusBadRequest, jsonResponse.Json{
				"message": "file not found",
			})
			return r
		}
		r.Error(err)
		return r
	}

	if d.File != nil {
		file, fileOpenErr := d.File.Open()
		if fileOpenErr != nil {
			r.Error(fileOpenErr)
			return r
		}
		defer file.Close()

		contentType, contentTypeErr := utils.GetFileContentType(file)
		if contentTypeErr != nil {
			r.Error(contentTypeErr)
			return r
		}
		removeCtx, removeCancel := utils.GetMinioContext()
		defer removeCancel()

		removeErr := minioClient.Conn.RemoveObject(removeCtx, minioClient.ImageBucket, f.FileId, minio.RemoveObjectOptions{})
		if removeErr != nil {
			r.Error(removeErr)
			return r
		}

		putCtx, putCancel := utils.GetMinioContext()
		defer putCancel()
		info, putErr := minioClient.Conn.PutObject(putCtx, minioClient.ImageBucket, f.FileId, file, d.File.Size, minio.PutObjectOptions{ContentType: contentType})
		if putErr != nil {
			r.Error(putErr)
			return r
		}

		dbCtx, dbCancel := utils.GetDatabaseContext()
		defer dbCancel()

		dbErr := db.Poll.QueryRow(dbCtx, `--sql
			UPDATE files SET hash_md5 = $1, content_type = $2 WHERE id = $3;
		`, info.ETag, contentType, d.Id).Scan()

		if dbErr != nil && dbErr != pgx.ErrNoRows {
			r.Error(dbErr)
			return r
		}

	}

	if d.Filename != nil {
		dbCtx, dbCancel := utils.GetDatabaseContext()
		defer dbCancel()

		dbErr := db.Poll.QueryRow(dbCtx, `--sql
			UPDATE files SET filename = $1 WHERE id = $2;
		`, d.Filename, d.Id).Scan()

		if dbErr != nil && dbErr != pgx.ErrNoRows {
			r.Error(dbErr)
			return r
		}
	}

	re := &emptyResponse.Response{}
	re.Write(http.StatusNoContent)
	return re
}

func DeleteFile(d dtos.DeleteFile) types.Response {
	r := &jsonResponse.Response{}

	f := models.File{}
	dbFileCtx, dbFileCancel := utils.GetDatabaseContext()
	defer dbFileCancel()
	dbFileErr := db.Poll.QueryRow(dbFileCtx, `--sql
		SELECT file_id FROM files WHERE id = $1;
	`, d.Id).Scan(&f.FileId)

	if dbFileErr != nil {
		if dbFileErr == pgx.ErrNoRows {
			r.Write(http.StatusBadRequest, jsonResponse.Json{
				"message": "file not found",
			})
			return r
		}
		r.Error(dbFileErr)
		return r
	}

	minioCtx, minioCancel := utils.GetMinioContext()
	defer minioCancel()
	minioErr := minioClient.Conn.RemoveObject(minioCtx, minioClient.ImageBucket, f.FileId, minio.RemoveObjectOptions{})

	if minioErr != nil {
		r.Error(minioErr)
		return r
	}

	dbCtx, dbCancel := utils.GetDatabaseContext()
	defer dbCancel()

	dbErr := db.Poll.QueryRow(dbCtx, `--sql
		DELETE FROM files WHERE id = $1;
	`, d.Id).Scan()

	if dbErr != nil && dbErr != pgx.ErrNoRows {
		r.Error(dbErr)
		return r
	}

	re := &emptyResponse.Response{}
	re.Write(http.StatusNoContent)
	return re
}
