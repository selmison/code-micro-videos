package rest

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/listing"
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func CategoryHandler(r *chi.Mux, l listing.Service, m modifying.Service) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", GetCategories(l))
		r.Post("/", CreateCategory(m))
		r.Route("/{categoryName}", func(r chi.Router) {
			r.Use(CategoryCtx(l))
			r.Get("/", GetCategory)
			r.Put("/", UpdateCategory(m))
			r.Delete("/", DeleteCategory(m))
		})
	})
}

func CreateCategory(s modifying.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryRequest := &CategoryRequest{}
		if err := render.Bind(r, categoryRequest); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		categoryDTO := categoryRequest.CategoryDTO
		if err := s.AddCategory(*categoryDTO); err != nil {
			if errors.Is(err, modifying.ErrAlreadyExists) {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		render.Status(r, http.StatusCreated)
		render.Render(w, r, NewCategoryDTOResponse(categoryDTO))
	}
}

func GetCategories(s listing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cs, err := s.GetCategories(math.MaxInt8)
		if err != nil {
			log.Fatalln(err)
		}
		if err := render.RenderList(w, r, NewCategoryListResponse(cs)); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
	}
}

func CategoryCtx(s listing.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var category models.Category
			var err error
			if categoryName := chi.URLParam(r, "categoryName"); strings.TrimSpace(categoryName) != "" {
				if r.Method == http.MethodPut {
					categoryRequest := &CategoryRequest{}
					if err := render.Bind(r, categoryRequest); err != nil {
						render.Render(w, r, ErrInvalidRequest(err))
						return
					}
					categoryDTO := *categoryRequest.CategoryDTO
					ctx := context.WithValue(r.Context(), "categoryName", categoryName)
					ctx = context.WithValue(ctx, "categoryDTO", categoryDTO)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				if r.Method == http.MethodDelete {
					ctx := context.WithValue(r.Context(), "categoryName", categoryName)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				category, err = s.FetchCategory(categoryName)
			} else {
				render.Render(w, r, ErrNotFound)
				return
			}
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}
			ctx := context.WithValue(r.Context(), "category", category)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	category := r.Context().Value("category").(models.Category)
	if err := render.Render(w, r, NewCategoryResponse(&category)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func UpdateCategory(s modifying.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryName := r.Context().Value("categoryName").(string)
		categoryDTO := r.Context().Value("categoryDTO").(modifying.CategoryDTO)
		if err := s.UpdateCategory(categoryName, categoryDTO); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
		render.Render(w, r, NewCategoryDTOResponse(&categoryDTO))
	}
}

func DeleteCategory(s modifying.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryName := r.Context().Value("categoryName").(string)
		if err := s.RemoveCategory(categoryName); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
	}
}
