package expense

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name *string `json:"name" validate:"required"`
}

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
}

func (CategoryResponse) FromEntity(category CategoryEntity) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		IsDefault: category.IsDefault,
	}
}
