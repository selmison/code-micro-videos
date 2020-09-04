package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s *server) handleCategoryCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryDTO := &crud.CategoryDTO{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.errInternalServer(w, err)
		}
		if err := r.Body.Close(); err != nil {
			s.errInternalServer(w, err)
		}
		if err := json.Unmarshal(body, &categoryDTO); err != nil {
			s.errUnprocessableEntity(w, err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				s.errInternalServer(w, err)
			}
		}
		if err := s.svc.AddCategory(*categoryDTO); err != nil {
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
		categories, err := s.svc.GetCategories(math.MaxInt8)
		if err != nil {
			s.errInternalServer(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		categoriesDTO := make([]crud.CategoryDTO, len(categories))
		for i, category := range categories {
			categoriesDTO[i] = crud.CategoryDTO{
				Name:        category.Name,
				Description: category.Description.String,
			}
		}
		if err := json.NewEncoder(w).Encode(categoriesDTO); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleCategoryGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category models.Category
		var err error
		params := httprouter.ParamsFromContext(r.Context())
		if categoryName := params.ByName("name"); strings.TrimSpace(categoryName) != "" {
			category, err = s.svc.FetchCategory(categoryName)
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
		categoryDTO := crud.CategoryDTO{
			Name:        category.Name,
			Description: category.Description.String,
		}
		if err := json.NewEncoder(w).Encode(categoryDTO); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleCategoryUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		categoryDTO := &crud.CategoryDTO{}
		if err := s.bodyToStruct(w, r, categoryDTO); err != nil {
			return
		}
		params := httprouter.ParamsFromContext(r.Context())
		if categoryName := params.ByName("name"); strings.TrimSpace(categoryName) != "" {
			err = s.svc.UpdateCategory(categoryName, *categoryDTO)
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
			err = s.svc.RemoveCategory(categoryName)
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
