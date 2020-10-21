package genre

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct {
	NewGenre NewGenre
}

type createResponse struct {
	Genre Genre `json:"genre"`
	Err   error `json:"error,omitempty"`
}

func (c createResponse) Error() error { return c.Err }

func MakeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		c, err := svc.Create(ctx, req.NewGenre)
		if err != nil {
			return createResponse{Genre{}, err}, nil
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
	Genres []Genre `json:"genres"`
	Err    error   `json:"error,omitempty"`
}

func (c listResponse) Error() error { return c.Err }

func MakeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		genres, err := svc.List(ctx)
		if err != nil {
			return listResponse{nil, err}, nil
		}
		return listResponse{genres, nil}, nil
	}
}

type ShowRequest struct {
	Id string
}

type showResponse struct {
	Genre Genre `json:"genre"`
	Err   error `json:"error,omitempty"`
}

func (c showResponse) Error() error { return c.Err }

func MakeShowEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShowRequest)
		genre, err := svc.Show(ctx, req.Id)
		if err != nil {
			return showResponse{Genre{}, err}, nil
		}
		return showResponse{genre, nil}, nil
	}
}

type UpdateRequest struct {
	Id          string
	UpdateGenre UpdateGenre
}

type updateResponse struct {
	Err error `json:"error,omitempty"`
}

func (c updateResponse) Error() error { return c.Err }

func MakeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := svc.Update(ctx, req.Id, req.UpdateGenre)
		if err != nil {
			return updateResponse{err}, nil
		}
		return updateResponse{nil}, nil
	}
}
