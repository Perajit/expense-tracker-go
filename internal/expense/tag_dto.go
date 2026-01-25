package expense

type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateTagRequest struct {
	Name *string `json:"name"`
}

type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (TagResponse) FromEntity(tag TagEntity) TagResponse {
	return TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}
