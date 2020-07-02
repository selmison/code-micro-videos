package rest

import (
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"net/http"
)

type CategoryRequest struct {
	*modifying.CategoryDTO
}

func NewCategoryRequest(c *modifying.CategoryDTO) *CategoryRequest {
	return &CategoryRequest{c}
}

func (c *CategoryRequest) Bind(r *http.Request) error {
	if err := c.CategoryDTO.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *CategoryRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//func NewCategoryListRequest(cs modifying.CategoryDTO) []render.Renderer {
//	list := make([]render.Renderer, len(cs))
//	for i, c := range cs {
//		list[i] = NewCategoryResponse(c)
//	}
//	return list
//}
