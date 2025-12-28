package dtos

type CreateProductForm struct {
	ID          int     `form:"product_id"`
	Name        string  `form:"product_name"`
	CategoryID  int     `form:"product_category"`
	Price       float64 `form:"product_price"`
	Stock       int     `form:"product_stock"`
	Description string  `form:"product_description"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
}

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	CategoryID  int      `json:"category_id"`
	Price       float64  `json:"price"`
	Images      []Image  `json:"images"`
	IsActive    bool     `json:"is_active"`
	SKU         string   `json:"sku"`
	Stock       int      `json:"stock"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
}

type Image struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductResponse struct {
	SharedResponse
	Product []Product `json:"data"`
	Meta    Meta      `json:"meta"`
}

type SingleProductResponse struct {
	SharedResponse
	Product Product `json:"data"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ProductImagesResponse struct {
	SharedResponse
	Images map[string]string `json:"data"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
}
