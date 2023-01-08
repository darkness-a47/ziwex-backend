package models

type File struct {
	Filename    string `json:"filename,omitempty"`
	FileId      string `json:"file_id,omitempty"`
	HashMd5     string `json:"hash_md5,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}
