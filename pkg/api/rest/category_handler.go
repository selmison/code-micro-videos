package rest

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/julienschmidt/httprouter"

	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s *server) handleCategoryCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := CategoryDTO{}
		if err := render.Decode(r, &dto); err != nil {
			s.errBadRequest(w, err)
			return
		}
		if err := s.svc.CreateCategory(r.Context(), *mapToSvcCategory(dto)); err != nil {
			if errors.Is(err, logger.ErrIsRequired) {
				s.errBadRequest(w, err)
				return
			}
			if errors.Is(err, logger.ErrAlreadyExists) {
				s.errStatusConflict(w, err)
				return
			}
			if errors.Is(err, logger.ErrNotFound) {
				s.errNotFound(w, err)
				return
			}
			s.errInternalServer(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(http.StatusText(http.StatusCreated)); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleCategoriesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := s.svc.GetCategories(r.Context(), math.MaxInt8)
		if err != nil {
			s.errInternalServer(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		categoriesDTO := make([]service.Category, len(categories))
		for i, category := range categories {
			categoriesDTO[i] = service.Category{
				Name:        category.Name,
				Description: category.Description,
			}
		}
		if err := json.NewEncoder(w).Encode(categoriesDTO); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleCategoryGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category service.Category
		var err error
		params := httprouter.ParamsFromContext(r.Context())
		if categoryName := params.ByName("name"); strings.TrimSpace(categoryName) != "" {
			category, err = s.svc.FetchCategory(r.Context(), categoryName)
			if err != nil {
				if errors.Is(err, logger.ErrNotFound) {
					s.errNotFound(w, err)
					return
				}
				if errors.Is(err, logger.ErrInternalApplication) {
					s.errInternalServer(w, err)
					return
				}
			}
		} else {
			s.errBadRequest(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		categoryDTO := service.Category{
			Name:        category.Name,
			Description: category.Description,
		}
		if err := json.NewEncoder(w).Encode(categoryDTO); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleCategoryUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		categoryDTO := &service.Category{}
		if err := s.bodyToStruct(w, r, categoryDTO); err != nil {
			return
		}
		params := httprouter.ParamsFromContext(r.Context())
		if categoryName := params.ByName("name"); strings.TrimSpace(categoryName) != "" {
			err = s.svc.UpdateCategory(r.Context(), categoryName, *categoryDTO)
			if err != nil {
				if errors.Is(err, logger.ErrNotFound) {
					s.errNotFound(w, err)
					return
				}
				if errors.Is(err, logger.ErrInternalApplication) {
					s.errInternalServer(w, err)
					return
				}
			}
		} else {
			s.errBadRequest(w, err)
			return
		}
	}
}

func (s *server) handleCategoryDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		params := httprouter.ParamsFromContext(r.Context())
		if categoryName := params.ByName("name"); strings.TrimSpace(categoryName) != "" {
			err = s.svc.RemoveCategory(r.Context(), categoryName)
			if err != nil {
				if errors.Is(err, logger.ErrNotFound) {
					s.errNotFound(w, err)
					return
				}
				if errors.Is(err, logger.ErrInternalApplication) {
					s.errInternalServer(w, err)
					return
				}
			}
		} else {
			s.errBadRequest(w, err)
			return
		}
	}
}

func mapToSvcCategory(dto CategoryDTO) *service.Category {
	var genres []service.GenreOfCategory
	if dto.Genres != nil {
		for _, genre := range *dto.Genres {
			genres = append(genres, service.GenreOfCategory(genre))
		}
	}
	return &service.Category{
		Description: dto.Description,
		Genres:      &genres,
		Id:          dto.Id,
		Name:        dto.Name,
	}
}
