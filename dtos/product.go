package dtos

type ProductKeyValue struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type CreateProduct struct {
	Url                 string            `json:"url" validate:"required"`
	Title               string            `json:"title" validate:"required"`
	Description         *string           `json:"description" validate:"required"`
	DescriptionKeyValue []ProductKeyValue `json:"description_key_value" validate:"required"`
	Options             []ProductKeyValue `json:"options" validate:"required"`
	Price               float64           `json:"price" validate:"required"`
	RelatedProducts     []int             `json:"related_products" validate:"required"`
	RecommendProducts   []int             `json:"recommend_products" validate:"required"`
	Images              []int             `json:"images" validate:"required"`
	MainImage           int               `json:"main_image" validate:"required"`
}
