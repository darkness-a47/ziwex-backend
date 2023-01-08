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

type UpdateFile struct {
	Id       int     `form:"id" validate:"required,number"`
	Filename *string `form:"filename"`
	File     *multipart.FileHeader
}

type DeleteFile struct {
	Id int `query:"id" validate:"required,number"`
}
