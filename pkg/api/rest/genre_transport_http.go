package rest

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/selmison/code-micro-videos/pkg/genre"
)

func NewGenreRoutes(router Router, svc genre.Service) {
	router.AddRoute("GET", "/genres",
		httptransport.NewServer(
			genre.MakeListEndpoint(svc),
			decodeGenreListRequest,
			encodeResponse,
		),
	)
	router.AddRoute("GET", "/genres/:id",
		httptransport.NewServer(
			genre.MakeShowEndpoint(svc),
			decodeGenreShowRequest,
			encodeResponse,
		),
	)
	router.AddRoute("POST", "/genres/",
		httptransport.NewServer(
			genre.MakeCreateEndpoint(svc),
			decodeGenreCreateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("PUT", "/genres/:id",
		httptransport.NewServer(
			genre.MakeUpdateEndpoint(svc),
			decodeGenreUpdateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("DELETE", "/genres/:id",
		httptransport.NewServer(
			genre.MakeDestroyEndpoint(svc),
			decodeGenreDestroyRequest,
			encodeResponse,
		),
	)
}

func decodeGenreCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request genre.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGenreDestroyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request genre.DestroyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGenreListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request genre.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGenreShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request genre.ShowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGenreUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request genre.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
