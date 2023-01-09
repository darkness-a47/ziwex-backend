package models

type ProductKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Product struct {
	Id                  int               `json:"id,omitempty"`
	Url                 string            `json:"url,omitempty"`
	Title               string            `json:"title,omitempty"`
	Description         string            `json:"description,omitempty"`
	DescriptionKeyValue []ProductKeyValue `json:"description_key_value,omitempty"`
	Options             []ProductKeyValue `json:"options,omitempty"`
	Price               float64           `json:"price,omitempty"`
	RelatedProducts     []Product         `json:"related_products,omitempty"`
	RecommendProducts   []Product         `json:"recommend_products,omitempty"`
	Categories          []Category        `json:"categories,omitempty"`
	Images              []File            `json:"images,omitempty"`
	MainImageIndex      int               `json:"main_image_index,omitempty"`
}
