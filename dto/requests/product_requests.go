package requests

type GetProductByIdRequest struct {
	ID string `json:"id" binding:"required,uuid"`
}

type GetProductBySlugRequest struct {
	Slug string `json:"slug" binding:"required,max=100"`
}

type GetProductsByCategoryRequest struct {
	CategoryID string `json:"category_id" binding:"required,uuid"`
}

type GetProductsBySubcategoryRequest struct {
	SubcategoryID string `json:"subcategory_id" binding:"required,uuid"`
}
