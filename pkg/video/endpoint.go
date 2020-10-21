package video

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct {
	NewVideo NewVideo
}

type createResponse struct {
	Video Video `json:"video"`
	Err   error `json:"error,omitempty"`
}

func (c createResponse) Error() error { return c.Err }

func MakeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		c, err := svc.Create(ctx, req.NewVideo)
		if err != nil {
			return createResponse{Video{}, err}, nil
		}
		return createResponse{c, nil}, nil
	}
}

type DestroyRequest struct {
	Id string
}

type destroyResponse struct {
	Err error `json:"error,omitempty"`
}

func (c destroyResponse) Error() error { return c.Err }

func MakeDestroyEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DestroyRequest)
		err := svc.Destroy(ctx, req.Id)
		if err != nil {
			return destroyResponse{err}, nil
		}
		return destroyResponse{nil}, nil
	}
}

type ListRequest struct{}

type listResponse struct {
	Videos []Video `json:"videos"`
	Err    error   `json:"error,omitempty"`
}

func (c listResponse) Error() error { return c.Err }

func MakeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		videos, err := svc.List(ctx)
		if err != nil {
			return listResponse{nil, err}, nil
		}
		return listResponse{videos, nil}, nil
	}
}

type ShowRequest struct {
	Id string
}

type showResponse struct {
	Video Video `json:"video"`
	Err   error `json:"error,omitempty"`
}

func (c showResponse) Error() error { return c.Err }

func MakeShowEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShowRequest)
		video, err := svc.Show(ctx, req.Id)
		if err != nil {
			return showResponse{Video{}, err}, nil
		}
		return showResponse{video, nil}, nil
	}
}

type UpdateRequest struct {
	Id          string
	UpdateVideo UpdateVideo
}

type updateResponse struct {
	Err error `json:"error,omitempty"`
}

func (c updateResponse) Error() error { return c.Err }

func MakeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := svc.Update(ctx, req.Id, req.UpdateVideo)
		if err != nil {
			return updateResponse{err}, nil
		}
		return updateResponse{nil}, nil
	}
}
