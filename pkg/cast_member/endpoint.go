package cast_member

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct {
	NewCastMemberDTO
}

type CreateResponse struct {
	CastMember CastMemberMap `json:"cast_member"`
	Err        error         `json:"error,omitempty"`
}

func (c CreateResponse) Error() error { return c.Err }

func MakeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		castMember, err := svc.Create(ctx, req.NewCastMemberDTO)
		if err != nil {
			return CreateResponse{nil, err}, nil
		}
		castMemberMap := CastMemberMap{
			"cast_member": castMember,
		}
		//castMemberDTO := castMemberDTO{
		//	Id:   castMember.Id(),
		//	Name: castMember.Name(),
		//	Type: castMember.Type(),
		//}
		return CreateResponse{castMemberMap, nil}, nil
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

type CastMemberMap map[string]CastMember

func (cm *CastMemberMap) UnmarshalJSON(data []byte) error {
	castMembers := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &castMembers)
	if err != nil {
		return err
	}
	result := make(CastMemberMap)
	for k, v := range castMembers {
		//fmt.Println(k, string(v))
		switch k {
		case "cast_member":
			castMember := castMember{}
			err := json.Unmarshal(v, &castMember)
			if err != nil {
				return err
			}
			//fmt.Println("UnmarshalJSON", castMember)
			result[k] = &castMember
		default:
			return errors.New("unrecognized shape")
		}
	}
	*cm = result
	return nil
}

type ListResponse struct {
	CastMembers []CastMemberMap `json:"cast_members"`
	Err         error           `json:"error,omitempty"`
}

func (c ListResponse) Error() error { return c.Err }

func MakeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		castMembers, err := svc.List(ctx)
		castMemberMaps := make([]CastMemberMap, len(castMembers))
		for i, castMember := range castMembers {
			castMemberMaps[i] = CastMemberMap{
				"cast_member": castMember,
			}
		}
		//CastMemberDTOs := make([]castMember, len(castMembers))
		//for i, castMember := range castMembers {
		//	CastMemberDTOs[i] = castMemberDTO{
		//		Id:   castMember.Id(),
		//		Name: castMember.Name(),
		//		Type: castMember.Type(),
		//	}
		//}
		if err != nil {
			return ListResponse{nil, err}, nil
		}
		return ListResponse{castMemberMaps, nil}, nil
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
