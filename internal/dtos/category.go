package dtos

type CategoryResponse struct {
	SharedResponse
	Categories []Category `json:"data"`
}

type SingleCategoryResponse struct {
	SharedResponse
	Category Category `json:"data"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryForm struct {
	Name        string `form:"category_name"`
	Description string `form:"category_description"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
