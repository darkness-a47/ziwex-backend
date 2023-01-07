package dtos

type CreateCategory struct {
	Title            string   `json:"title" validate:"required"`
	ImageUrl         string   `json:"image_url" validate:"required"`
	Description      string   `json:"description" validate:"required"`
	ParentCategoryId *int     `json:"parent_category_id"`
	Tags             []string `json:"tags" validate:"required"`
}

type GetCategories struct {
	ParentCategoryId *int `query:"parent_category_id"`
	Page             int  `query:"page" validate:"required"`
	DataPerPage      int  `query:"data_per_page" validate:"required"`
}
