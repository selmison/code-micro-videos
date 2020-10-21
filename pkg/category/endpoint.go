package category

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct {
	NewCategory NewCategory
}

type createResponse struct {
	Category Category `json:"category"`
	Err      error    `json:"error,omitempty"`
}

func (c createResponse) Error() error { return c.Err }

func MakeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		c, err := svc.Create(ctx, req.NewCategory)
		if err != nil {
			return createResponse{Category{}, err}, nil
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
	Categories []Category `json:"categories"`
	Err        error      `json:"error,omitempty"`
}

func (c listResponse) Error() error { return c.Err }

func MakeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		categories, err := svc.List(ctx)
		if err != nil {
			return listResponse{nil, err}, nil
		}
		return listResponse{categories, nil}, nil
	}
}

type ShowRequest struct {
	Id string
}

type showResponse struct {
	Category Category `json:"category"`
	Err      error    `json:"error,omitempty"`
}

func (c showResponse) Error() error { return c.Err }

func MakeShowEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShowRequest)
		category, err := svc.Show(ctx, req.Id)
		if err != nil {
			return showResponse{Category{}, err}, nil
		}
		return showResponse{category, nil}, nil
	}
}

type UpdateRequest struct {
	Id             string
	UpdateCategory UpdateCategory
}

type updateResponse struct {
	Err error `json:"error,omitempty"`
}

func (c updateResponse) Error() error { return c.Err }

func MakeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := svc.Update(ctx, req.Id, req.UpdateCategory)
		if err != nil {
			return updateResponse{err}, nil
		}
		return updateResponse{nil}, nil
	}
}
