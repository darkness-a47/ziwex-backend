package models

type Category struct {
	Id               int      `json:"id,omitempty"`
	Title            string   `json:"title,omitempty"`
	ImageId          string   `json:"image_id,omitempty"`
	Description      string   `json:"description,omitempty"`
	ParentCategoryId *int     `json:"parent_category_id,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}
