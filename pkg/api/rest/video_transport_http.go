package rest

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/selmison/code-micro-videos/pkg/video"
)

const (
	MaxMemory      = 10 << 20
	VideoFileField = "video_file"
)

func NewVideoRoutes(router Router, svc video.Service) {
	router.AddRoute("GET", "/videos",
		httptransport.NewServer(
			video.MakeListEndpoint(svc),
			decodeVideoListRequest,
			encodeResponse,
		),
	)
	router.AddRoute("GET", "/videos/:id",
		httptransport.NewServer(
			video.MakeShowEndpoint(svc),
			decodeVideoShowRequest,
			encodeResponse,
		),
	)
	router.AddRoute("POST", "/videos/",
		httptransport.NewServer(
			video.MakeCreateEndpoint(svc),
			decodeVideoCreateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("PUT", "/videos/:id",
		httptransport.NewServer(
			video.MakeUpdateEndpoint(svc),
			decodeVideoUpdateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("DELETE", "/videos/:id",
		httptransport.NewServer(
			video.MakeDestroyEndpoint(svc),
			decodeVideoDestroyRequest,
			encodeResponse,
		),
	)
}

func decodeVideoCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request video.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVideoDestroyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request video.DestroyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVideoListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request video.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVideoShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request video.ShowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVideoUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request video.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
