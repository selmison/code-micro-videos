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

func (s *server) handleVideoCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoDTO := &crud.VideoDTO{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.errInternalServer(w, err)
		}
		if err := r.Body.Close(); err != nil {
			s.errInternalServer(w, err)
		}
		if err := json.Unmarshal(body, &videoDTO); err != nil {
			s.errUnprocessableEntity(w, err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				s.errInternalServer(w, err)
			}
		}
		if err := s.svc.AddVideo(*videoDTO); err != nil {
			if errors.Is(err, logger.ErrIsRequired) {
				s.errBadRequest(w, err)
				return
			}
			if errors.Is(err, logger.ErrAlreadyExists) {
				s.errStatusConflict(w, err)
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

func (s *server) handleVideosGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videos, err := s.svc.GetVideos(math.MaxInt8)
		if err != nil {
			s.errInternalServer(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		videosDTO := make([]crud.VideoDTO, len(videos))
		for i, video := range videos {
			videosDTO[i] = crud.VideoDTO{
				Title:        video.Title,
				Description:  video.Description,
				YearLaunched: video.YearLaunched,
				Opened:       video.Opened.Bool,
				Rating:       crud.VideoRating(video.Rating),
				Duration:     video.Duration,
			}
		}
		if err := json.NewEncoder(w).Encode(videosDTO); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleVideoGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var video models.Video
		var err error
		params := httprouter.ParamsFromContext(r.Context())
		if videoTitle := params.ByName("title"); strings.TrimSpace(videoTitle) != "" {
			video, err = s.svc.FetchVideo(videoTitle)
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
		videoDTO := crud.VideoDTO{
			Title:        video.Title,
			Description:  video.Description,
			YearLaunched: video.YearLaunched,
			Opened:       video.Opened.Bool,
			Rating:       crud.VideoRating(video.Rating),
			Duration:     video.Duration,
		}
		if err := json.NewEncoder(w).Encode(videoDTO); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) handleVideoUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		videoDTO := &crud.VideoDTO{}
		if err := s.bodyToStruct(w, r, videoDTO); err != nil {
			return
		}
		params := httprouter.ParamsFromContext(r.Context())
		if videoTitle := params.ByName("title"); strings.TrimSpace(videoTitle) != "" {
			err = s.svc.UpdateVideo(videoTitle, *videoDTO)
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

func (s *server) handleVideoDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		params := httprouter.ParamsFromContext(r.Context())
		if videoTitle := params.ByName("title"); strings.TrimSpace(videoTitle) != "" {
			err = s.svc.RemoveVideo(videoTitle)
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
