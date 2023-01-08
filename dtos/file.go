package dtos

import "mime/multipart"

type UploadFile struct {
	Filename string `form:"filename" validate:"required"`
	File     *multipart.FileHeader
}

type ServeFile struct {
	FileId string `param:"file_id" validate:"required,uuid"`
}

type GetFiles struct {
	Filename    *string `query:"filename" validate:"required"`
	Page        int     `query:"page" validate:"required"`
	DataPerPage int     `query:"data_per_page" validate:"required"`
}
