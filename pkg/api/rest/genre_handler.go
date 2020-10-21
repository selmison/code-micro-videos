package rest

//import (
//	"io/ioutil"
//	"net/http"
//	"strings"
//
//	"github.com/julienschmidt/httprouter"
//
//	"github.com/selmison/code-micro-videos/pkg/logger"
//)
//
//func (s *server) handleGenreCreate() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		genreDTO := &service.Genre{}
//		body, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			s.errInternalServer(w, err)
//		}
//		if err := r.Body.Close(); err != nil {
//			s.errInternalServer(w, err)
//		}
//		if err := json.Unmarshal(body, &genreDTO); err != nil {
//			s.errUnprocessableEntity(w, err)
//			if err := json.NewEncoder(w).Encode(err); err != nil {
//				s.errInternalServer(w, err)
//			}
//		}
//		if err := s.svc.CreateGenre(r.Context(), *genreDTO); err != nil {
//			if errors.Is(err, logger.ErrIsRequired) {
//				s.errBadRequest(w, err)
//				return
//			}
//			if errors.Is(err, logger.ErrAlreadyExists) {
//				s.errStatusConflict(w, err)
//				return
//			}
//			s.errInternalServer(w, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusCreated)
//		if err := json.NewEncoder(w).Encode(http.StatusText(http.StatusCreated)); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleGenresGet() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		genres, err := s.svc.GetGenres(nil, math.MaxInt8)
//		if err != nil {
//			s.errInternalServer(w, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//		genresDTO := make([]service.Genre, len(genres))
//		for i, genre := range genres {
//			genresDTO[i] = service.Genre{
//				Name: genre.Name,
//			}
//		}
//		if err := json.NewEncoder(w).Encode(genresDTO); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleGenreGet() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var genre service.Genre
//		var err error
//		params := httprouter.ParamsFromContext(r.Context())
//		if genreName := params.ByName("name"); strings.TrimSpace(genreName) != "" {
//			genre, err = s.svc.FetchGenre(r.Context(), genreName)
//			if err != nil {
//				if errors.Is(err, logger.ErrNotFound) {
//					s.errNotFound(w, err)
//					return
//				}
//				if errors.Is(err, logger.ErrInternalApplication) {
//					s.errInternalServer(w, err)
//					return
//				}
//			}
//		} else {
//			s.errBadRequest(w, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//		genreDTO := service.Genre{
//			Name: genre.Name,
//		}
//		if err := json.NewEncoder(w).Encode(genreDTO); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleGenreUpdate() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var err error
//		genreDTO := &service.Genre{}
//		if err := s.bodyToStruct(w, r, genreDTO); err != nil {
//			return
//		}
//		params := httprouter.ParamsFromContext(r.Context())
//		if genreName := params.ByName("name"); strings.TrimSpace(genreName) != "" {
//			err = s.svc.UpdateGenre(nil, genreName, *genreDTO)
//			if err != nil {
//				if errors.Is(err, logger.ErrNotFound) {
//					s.errNotFound(w, err)
//					return
//				}
//				if errors.Is(err, logger.ErrInternalApplication) {
//					s.errInternalServer(w, err)
//					return
//				}
//			}
//		} else {
//			s.errBadRequest(w, err)
//			return
//		}
//	}
//}
//
//func (s *server) handleGenreDelete() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var err error
//		params := httprouter.ParamsFromContext(r.Context())
//		if genreName := params.ByName("name"); strings.TrimSpace(genreName) != "" {
//			err = s.svc.RemoveGenre(nil, genreName)
//			if err != nil {
//				if errors.Is(err, logger.ErrNotFound) {
//					s.errNotFound(w, err)
//					return
//				}
//				if errors.Is(err, logger.ErrInternalApplication) {
//					s.errInternalServer(w, err)
//					return
//				}
//			}
//		} else {
//			s.errBadRequest(w, err)
//			return
//		}
//	}
//}
