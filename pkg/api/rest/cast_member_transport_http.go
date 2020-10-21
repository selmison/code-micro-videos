package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func NewCastMemberRoutes(router Router, svc cast_member.Service) {
	router.AddRoute("GET", "/cast_members",
		httptransport.NewServer(
			cast_member.MakeListEndpoint(svc),
			decodeCastMemberListRequest,
			encodeResponse,
		),
	)
	router.AddRoute("GET", "/cast_members/:id",
		httptransport.NewServer(
			cast_member.MakeShowEndpoint(svc),
			decodeCastMemberShowRequest,
			encodeResponse,
		),
	)
	router.AddRoute("POST", "/cast_members/",
		httptransport.NewServer(
			cast_member.MakeCreateEndpoint(svc),
			decodeCastMemberCreateRequest,
			encodeCastMemberCreateResponse,
		),
	)
	router.AddRoute("PUT", "/cast_members/:id",
		httptransport.NewServer(
			cast_member.MakeUpdateEndpoint(svc),
			decodeCastMemberUpdateRequest,
			encodeResponse,
		),
	)
	router.AddRoute("DELETE", "/cast_members/:id",
		httptransport.NewServer(
			cast_member.MakeDestroyEndpoint(svc),
			decodeCastMemberDestroyRequest,
			encodeResponse,
		),
	)
}

func decodeCastMemberCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request cast_member.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeCastMemberCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(logger.Errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

func decodeCastMemberDestroyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request cast_member.DestroyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCastMemberListRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCastMemberShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request cast_member.ShowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeCastMemberListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(logger.Errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeCastMemberUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request cast_member.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		log.Panic("encodeError with nil error")
	}
	if e, ok := err.(*logger.ResultError); ok {
		switch e.Err {
		case logger.ErrNotFound:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		case
			logger.ErrAlreadyExists,
			logger.ErrCouldNotBeEmpty,
			logger.ErrIsRequired,
			logger.ErrIsNotValidated,
			logger.ErrInvalidedLimit:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
