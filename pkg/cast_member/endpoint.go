package cast_member

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct {
	NewCastMemberDTO
}

type CreateResponse struct {
	CastMember CastMemberDTO `json:"cast_member"`
	Err        error         `json:"error,omitempty"`
}

func (c CreateResponse) Error() error { return c.Err }

func MakeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		castMember, err := svc.Create(ctx, req.NewCastMemberDTO)
		if err != nil {
			return CreateResponse{CastMemberDTO{}, err}, nil
		}
		castMemberDTO := CastMemberDTO{
			Id:   castMember.Id(),
			Name: castMember.Name(),
			Type: castMember.Type(),
		}
		return CreateResponse{castMemberDTO, nil}, nil
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

type ListResponse struct {
	CastMembers []CastMemberDTO `json:"cast_members"`
	Err         error           `json:"error,omitempty"`
}

func (c ListResponse) Error() error { return c.Err }

func MakeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		castMembers, err := svc.List(ctx)
		castMemberDTOs := make([]CastMemberDTO, len(castMembers))
		for i, castMember := range castMembers {
			castMemberDTOs[i] = CastMemberDTO{
				Id:   castMember.Id(),
				Name: castMember.Name(),
				Type: castMember.Type(),
			}
		}
		if err != nil {
			return ListResponse{nil, err}, nil
		}
		return ListResponse{castMemberDTOs, nil}, nil
	}
}

type ShowRequest struct {
	Id string
}

type showResponse struct {
	CastMember CastMember `json:"cast_member"`
	Err        error      `json:"error,omitempty"`
}

func (c showResponse) Error() error { return c.Err }

func MakeShowEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShowRequest)
		castMember, err := svc.Show(ctx, req.Id)
		if err != nil {
			return showResponse{nil, err}, nil
		}
		return showResponse{castMember, nil}, nil
	}
}

type UpdateRequest struct {
	Id               string
	UpdateCastMember UpdateCastMemberDTO
}

type updateResponse struct {
	Err error `json:"error,omitempty"`
}

func (c updateResponse) Error() error { return c.Err }

func MakeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := svc.Update(ctx, req.Id, req.UpdateCastMember)
		if err != nil {
			return updateResponse{err}, nil
		}
		return updateResponse{nil}, nil
	}
}
