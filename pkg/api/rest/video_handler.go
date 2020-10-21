package rest

//import (
//	"encoding/json"
//	"errors"
//	"math"
//	"net/http"
//	"strings"
//
//	"github.com/gorilla/schema"
//	"github.com/julienschmidt/httprouter"
//
//	"github.com/selmison/code-micro-videos/models"
//	"github.com/selmison/code-micro-videos/pkg/crud/service"
//	"github.com/selmison/code-micro-videos/pkg/logger"
//)
//
//var decoder = schema.NewDecoder()
//
//func (s *server) handleVideoCreate() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if err := r.ParseMultipartForm(MaxMemory); err != nil {
//			s.errInternalServer(w, err)
//			return
//		}
//		videoDTO := &service.VideoDTO{}
//		if err := decoder.Decode(videoDTO, r.PostForm); err != nil {
//			s.errInternalServer(w, err)
//		}
//		if err := r.Body.Close(); err != nil {
//			s.errInternalServer(w, err)
//		}
//		_, videoFile, err := r.FormFile(VideoFileField)
//		if err != nil {
//			s.errInternalServer(w, err)
//		}
//		videoDTO.File = videoFile
//		if _, err := s.svc.CreateVideo(ctx, *videoDTO); err != nil {
//			if errors.Is(err, logger.ErrIsRequired) {
//				s.errBadRequest(w, err)
//				return
//			}
//			if errors.Is(err, logger.ErrAlreadyExists) {
//				s.errStatusConflict(w, err)
//				return
//			}
//			if errors.Is(err, logger.ErrNotFound) {
//				s.errNotFound(w, err)
//				return
//			}
//			s.errInternalServer(w, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusCreated)
//		if _, err := w.Write([]byte(http.StatusText(http.StatusCreated))); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleVideosGet() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		videos, err := s.svc.GetVideos(math.MaxInt8)
//		if err != nil {
//			s.errInternalServer(w, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//		videosDTO := make([]*service.VideoDTO, len(videos))
//		for i, video := range videos {
//			dto, err := service.MapVideoToDTO(*video)
//			if err != nil {
//				s.errBadRequest(w, err)
//			}
//			videosDTO[i] = dto
//		}
//		if err := json.NewEncoder(w).Encode(videosDTO); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleVideoGet() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var video models.Video
//		var err error
//		params := httprouter.ParamsFromContext(r.Context())
//		if videoTitle := params.ByName("title"); strings.TrimSpace(videoTitle) != "" {
//			video, err = s.svc.FetchVideo(videoTitle)
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
//		videoDTO, err := service.MapVideoToDTO(video)
//		if err != nil {
//			s.errBadRequest(w, err)
//		}
//		if err := json.NewEncoder(w).Encode(videoDTO); err != nil {
//			s.errInternalServer(w, err)
//		}
//	}
//}
//
//func (s *server) handleVideoUpdate() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var err error
//		videoDTO := &service.VideoDTO{}
//		if err := s.bodyToStruct(w, r, videoDTO); err != nil {
//			return
//		}
//		params := httprouter.ParamsFromContext(r.Context())
//		videoTitle := params.ByName("title")
//		_, err = s.svc.UpdateVideo(videoTitle, *videoDTO)
//		if err != nil {
//			if errors.Is(err, logger.ErrNotFound) {
//				s.errNotFound(w, err)
//				return
//			}
//			if errors.Is(err, logger.ErrInternalApplication) {
//				s.errInternalServer(w, err)
//				return
//			}
//			s.errBadRequest(w, err)
//			return
//		}
//	}
//}
//
//func (s *server) handleVideoDelete() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var err error
//		params := httprouter.ParamsFromContext(r.Context())
//		if videoTitle := params.ByName("title"); strings.TrimSpace(videoTitle) != "" {
//			err = s.svc.RemoveVideo(videoTitle)
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
