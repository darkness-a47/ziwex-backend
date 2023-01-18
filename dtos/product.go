package dtos

type ProductKeyValue struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type CreateProduct struct {
	Url                 string            `json:"url" validate:"required"`
	Title               string            `json:"title" validate:"required"`
	Description         *string           `json:"description" validate:"required"`
	DescriptionKeyValue []ProductKeyValue `json:"description_key_value" validate:"required,dive"`
	Options             []ProductKeyValue `json:"options" validate:"required,dive"`
	Price               float64           `json:"price" validate:"required"`
	RelatedProducts     []int             `json:"related_products" validate:"required"`
	RecommendProducts   []int             `json:"recommend_products" validate:"required"`
	Categories          []int             `json:"categories" validate:"required"`
	Images              []int             `json:"images" validate:"required,gte=1"`
	MainImageIndex      *int              `json:"main_image_index" validate:"required,min=0"`
}

type GetProductsSummery struct {
	CategoryId  *int `query:"category_id"`
	Page        int  `query:"page" validate:"required"`
	DataPerPage int  `query:"data_per_page" validate:"required"`
}

type GetProductData struct {
	ProductUrl string `param:"product_url" validate:"required"`
}
