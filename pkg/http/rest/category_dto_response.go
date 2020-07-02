package rest

import (
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"net/http"
)

type CategoryDTOResponse struct {
	*modifying.CategoryDTO
}

func NewCategoryDTOResponse(c *modifying.CategoryDTO) *CategoryDTOResponse {
	return &CategoryDTOResponse{c}
}

func (c *CategoryDTOResponse) Bind(r *http.Request) error {
	return nil
}

func (c *CategoryDTOResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
