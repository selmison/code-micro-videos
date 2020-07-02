package rest

import (
	"github.com/go-chi/render"
	"github.com/selmison/code-micro-videos/models"
	"net/http"
)

type CategoryResponse struct {
	*models.Category
}

func NewCategoryResponse(c *models.Category) *CategoryResponse {
	return &CategoryResponse{c}
}

func (c *CategoryResponse) Bind(r *http.Request) error {
	return nil
}

func (c *CategoryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewCategoryListResponse(cs models.CategorySlice) []render.Renderer {
	list := make([]render.Renderer, len(cs))
	for i, c := range cs {
		list[i] = NewCategoryResponse(c)
	}
	return list
}
