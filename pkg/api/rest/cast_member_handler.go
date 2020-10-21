package rest

//import (
//	"encoding/json"
//	"errors"
//	"io/ioutil"
//	"math"
//	"net/http"
//	"strings"
//
//	"github.com/julienschmidt/httprouter"
//
//	"github.com/selmison/code-micro-videos/models"
//	"github.com/selmison/code-micro-videos/pkg/crud/service"
//	"github.com/selmison/code-micro-videos/pkg/logger"
//)
//
//func (s *server) handleCastMemberCreate() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		castMemberDTO := &service.CastMemberDTO{}
//		body, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			s.errInternalServer(w, err)
//		}
//		if err := r.Body.Close(); err != nil {
//			s.errInternalServer(w, err)
//		}
//		if err := json.Unmarshal(body, &castMemberDTO); err != nil {
//			s.errUnprocessableEntity(w, err)
//			if err := json.NewEncoder(w).Encode(err); err != nil {
//				s.errInternalServer(w, err)
//			}
//		}
//		if err := s.svc.AddCastMember(*castMemberDTO); err != nil {
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
//func (s *server) handleCastMembersGet() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		castMembers, err := s.svc.GetCastMembers(math.MaxInt8)
//		if err != nil {
//			s.errInternalServer(w, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//		castMembersDTO := make([]service.CastMemberDTO, len(castMembers))
//		for i, castMember := range castMembers {
//			castMembersDTO[i] = service.CastMemberDTO{
//				Name: castMember.Name,
//				Type: service.CastMemberType(castMember.Type),
//			}
//		}
//		if err := json.NewEncoder(w).Encode(castMembersDTO); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleCastMemberGet() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var castMember models.CastMember
//		var err error
//		params := httprouter.ParamsFromContext(r.Context())
//		if castMemberName := params.ByName("name"); strings.TrimSpace(castMemberName) != "" {
//			castMember, err = s.svc.FetchCastMember(castMemberName)
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
//		castMemberDTO := service.CastMemberDTO{
//			Name: castMember.Name,
//		}
//		if err := json.NewEncoder(w).Encode(castMemberDTO); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleCastMemberUpdate() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var err error
//		castMemberDTO := &service.CastMemberDTO{}
//		if err := s.bodyToStruct(w, r, castMemberDTO); err != nil {
//			return
//		}
//		params := httprouter.ParamsFromContext(r.Context())
//		if castMemberName := params.ByName("name"); strings.TrimSpace(castMemberName) != "" {
//			err = s.svc.UpdateCastMember(castMemberName, *castMemberDTO)
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
//func (s *server) handleCastMemberDelete() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var err error
//		params := httprouter.ParamsFromContext(r.Context())
//		if castMemberName := params.ByName("name"); strings.TrimSpace(castMemberName) != "" {
//			err = s.svc.RemoveCastMember(castMemberName)
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
