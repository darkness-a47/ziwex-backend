package models

type Category struct {
	Id               int      `json:"id,omitempty"`
	Title            string   `json:"title,omitempty"`
	ImageUrl         string   `json:"image_url,omitempty"`
	Description      string   `json:"description,omitempty"`
	ParentCategoryId *int     `json:"parent_category_id,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}
