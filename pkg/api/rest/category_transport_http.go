package rest

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/selmison/code-micro-videos/pkg/category"
)

func NewCategoryRoutes(router Router, svc category.Service) {
	router.AddRoute("GET", "/categories",
		httptransport.NewServer(
			category.MakeListEndpoint(svc),
			decodeCategoryListRequest,
			encodeResponse,
		),
	)
	router.AddRoute("GET", "/categories/:id",
		httptransport.NewServer(
			category.MakeShowEndpoint(svc),
			decodeCategoryShowRequest,
			encodeResponse,
		),
	)
	router.AddRoute("POST", "/categories/",
		httptransport.NewServer(
			category.MakeCreateEndpoint(svc),
			decodeCategoryCreateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("PUT", "/categories/:id",
		httptransport.NewServer(
			category.MakeUpdateEndpoint(svc),
			decodeCategoryUpdateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("DELETE", "/categories/:id",
		httptransport.NewServer(
			category.MakeDestroyEndpoint(svc),
			decodeCategoryDestroyRequest,
			encodeResponse,
		),
	)
}

func decodeCategoryCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request category.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCategoryDestroyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request category.DestroyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCategoryListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request category.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCategoryShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request category.ShowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCategoryUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request category.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
