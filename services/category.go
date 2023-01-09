package services

import (
	"net/http"
	"ziwex/db"
	"ziwex/dtos"
	"ziwex/models"
	"ziwex/types/jsonResponse"
	"ziwex/utils"

	"github.com/jackc/pgx/v5"
)

func CreateCategory(d dtos.CreateCategory) jsonResponse.Response {
	r := jsonResponse.Response{}

	fileCtx, fileCancel := utils.GetDatabaseContext()
	defer fileCancel()

	file := models.File{}
	fileErr := db.Poll.QueryRow(fileCtx, `--sql
		SELECT file_id FROM files where id = $1;
	`, d.ImageId).Scan(&file.FileId)

	if fileErr != nil {
		if fileErr == pgx.ErrNoRows {
			r.Write(http.StatusBadRequest, jsonResponse.Json{
				"message": "image not found",
			})
			return r
		}
		r.Error(fileErr)
		return r
	}

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	cat := models.Category{}
	err := db.Poll.QueryRow(ctx, `--sql
		INSERT INTO categories (title, image_id, description, parent_category_id, tags) VALUES ($1, $2, $3, $4, $5) RETURNING id;
		`, d.Title, file.FileId, d.Description, d.ParentCategoryId, d.Tags).Scan(&cat.Id)
	if err != nil {
		r.Error(err)
		return r
	}

	r.Write(http.StatusCreated, jsonResponse.Json{
		"message":     "category created",
		"category_id": cat.Id,
	})

	return r
}

func GetCategories(d dtos.GetCategories) jsonResponse.Response {
	res := jsonResponse.Response{}

	categories := make([]models.Category, 0)

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	offset := (d.Page - 1) * d.DataPerPage
	var rows pgx.Rows
	var err error
	if d.ParentCategoryId != nil {
		rows, err = db.Poll.Query(ctx, `--sql
			SELECT id, title, image_id, description, parent_category_id, tags, COUNT(*) OVER() AS total_count 
			FROM categories WHERE parent_category_id = $1 OFFSET $2 LIMIT $3
		`, d.ParentCategoryId, offset, d.DataPerPage)
	} else {
		rows, err = db.Poll.Query(ctx, `--sql
			SELECT id, title, image_id, description, parent_category_id, tags,  COUNT(*) OVER() AS total_count 
			FROM categories WHERE parent_category_id IS NULL OFFSET $1 LIMIT $2
		`, offset, d.DataPerPage)
	}
	if err != nil {
		res.Error(err)
		return res
	}

	var totalRows int
	for rows.Next() {
		c := models.Category{}
		err := rows.Scan(&c.Id, &c.Title, &c.ImageId, &c.Description, &c.ParentCategoryId, &c.Tags, &totalRows)
		if err != nil {
			res.Error(err)
			return res
		}
		categories = append(categories, c)
	}

	res.Write(http.StatusOK, jsonResponse.Json{
		"message":    "ok",
		"categories": categories,
		"total_rows": totalRows,
	})
	return res
}
