package cast_member

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/selmison/code-micro-videos/pkg/id_generator"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/validator"
)

type Service interface {
	// Create creates a new castMember.
	Create(ctx context.Context, newCastMember NewCastMemberDTO) (CastMember, error)

	// Destroy destroys a castMember.
	Destroy(ctx context.Context, name string) error

	// List returns a list of castMembers.
	List(ctx context.Context) ([]CastMember, error)

	// Show returns the details of a castMember.
	Show(ctx context.Context, id string) (CastMember, error)

	// Update updates an existing castMember.
	Update(ctx context.Context, id string, updateCastMember UpdateCastMemberDTO) error
}

type CastMemberType int8

const (
	Director CastMemberType = iota + 1
	Actor
)

func (c *CastMemberType) String() string {
	return [...]string{"Director", "Actor"}[*c-1]
}

func (c *CastMemberType) validate() error {
	switch *c {
	case Director, Actor:
		return nil
	}
	return fmt.Errorf("cast member type %w", logger.ErrIsNotValidated)
}

// CastMember represents a single CastMember.
type CastMember interface {
	Id() string
	Name() string
	Type() CastMemberType
	Update(updateCastMember UpdateCastMemberDTO)
	validate() error
	json.Marshaler
	json.Unmarshaler
}

type castMember struct {
	id             string
	name           string
	castMemberType CastMemberType
}

func NewCastMember(id string, newCastMember NewCastMemberDTO) (CastMember, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty)
	}
	if newCastMember.isEmpty() {
		return nil, fmt.Errorf("'%s' param %w", "newCastMember", logger.ErrCouldNotBeEmpty)
	}
	castMember := &castMember{
		id:             id,
		name:           strings.TrimSpace(newCastMember.Name),
		castMemberType: newCastMember.Type,
	}
	if err := castMember.validate(); err != nil {
		return nil, err
	}
	return castMember, nil
}

func (c *castMember) Id() string {
	return c.id
}

func (c *castMember) Name() string {
	return c.name
}

func (c *castMember) Type() CastMemberType {
	return c.castMemberType
}

func (c *castMember) UnmarshalJSON(b []byte) error {
	castMemberDTO := castMemberDTO{}
	if err := json.Unmarshal(b, &castMemberDTO); err != nil {
		return err
	}
	c.id = castMemberDTO.Id
	c.name = castMemberDTO.Name
	c.castMemberType = castMemberDTO.Type
	return nil
}

func (c castMember) MarshalJSON() ([]byte, error) {
	castMemberDTO := castMemberDTO{
		Id:   c.Id(),
		Name: c.Name(),
		Type: c.Type(),
	}
	return json.Marshal(castMemberDTO)
}

func (c *castMember) Update(updateCastMember UpdateCastMemberDTO) {
	if updateCastMember.Name != nil {
		c.name = *updateCastMember.Name
	}
	if updateCastMember.Type != nil {
		c.castMemberType = *updateCastMember.Type
	}
}

func (c *castMember) validate() error {
	if err := validator.CheckAll(
		validator.NewCheck(
			c.name == "",
			&logger.ResultError{
				ErrMsg: "the Name field",
				Err:    logger.ErrCouldNotBeEmpty,
			}),
	); err != nil {
		return err
	}
	if err := c.castMemberType.validate(); err != nil {
		return err
	}
	return nil
}

// castMemberDTO contains the details of a CastMember.
type castMemberDTO struct {
	Id   string         `json:"id"`
	Name string         `json:"name"`
	Type CastMemberType `json:"type"`
}

// NewCastMemberDTO contains the details of a new CastMember.
type NewCastMemberDTO struct {
	Name string         `json:"name"`
	Type CastMemberType `json:"type"`
}

// UpdateCastMemberDTO contains updates of an existing castMember.
type UpdateCastMemberDTO struct {
	Name *string         `json:"name"`
	Type *CastMemberType `json:"type"`
}

func (c *UpdateCastMemberDTO) validate() error {
	if c.Name != nil && strings.TrimSpace(*c.Name) == "" {
		return fmt.Errorf("the %s field %w", "Name", logger.ErrCouldNotBeEmpty)
	}
	if c.Type != nil {
		if err := c.Type.validate(); err != nil {
			return err
		}
	}
	return nil
}

type service struct {
	idGenerator id_generator.IdGenerator
	repo        Repository
}

// NewService returns a new Service with all of the expected middlewares wired in.
func NewService(idGenerator id_generator.IdGenerator, r Repository, logger log.Logger) Service {
	var svc Service
	{
		svc = service{idGenerator: idGenerator, repo: r}
		svc = NewValidationMiddleware()(svc)
		svc = NewLoggingMiddleware(logger)(svc)
	}
	return svc
}

func (svc service) Create(ctx context.Context, newCastMember NewCastMemberDTO) (CastMember, error) {
	id, err := svc.idGenerator.Generate()
	if err != nil {
		return nil, err
	}
	toCastMember, err := NewCastMember(id, newCastMember)
	if err != nil {
		return nil, err
	}
	err = svc.repo.Store(ctx, toCastMember)
	if err != nil {
		return nil, err
	}
	return toCastMember, nil
}

func (svc service) Destroy(ctx context.Context, id string) error {
	return svc.repo.DeleteOne(ctx, id)
}

func (svc service) List(ctx context.Context) ([]CastMember, error) {
	castMembers, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return castMembers, nil
}

func (svc service) Show(ctx context.Context, id string) (CastMember, error) {
	getCastMember, err := svc.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return getCastMember, nil
}

func (svc service) Update(ctx context.Context, id string, updateCastMember UpdateCastMemberDTO) error {
	err := svc.repo.UpdateOne(ctx, id, updateCastMember)
	if err != nil {
		return err
	}
	return nil
}

func (c NewCastMemberDTO) isEmpty() bool {
	return c.compare(NewCastMemberDTO{})
}

func (c NewCastMemberDTO) compare(b NewCastMemberDTO) bool {
	if &c == &b {
		return true
	}
	if c.Name != b.Name {
		return false
	}
	if c.Type != b.Type {
		return false
	}
	return true
}

func (c UpdateCastMemberDTO) isEmpty() bool {
	return c.compare(UpdateCastMemberDTO{})
}

func (c UpdateCastMemberDTO) compare(b UpdateCastMemberDTO) bool {
	if &c == &b {
		return true
	}
	if c.Name != b.Name {
		return false
	}
	if c.Type != b.Type {
		return false
	}
	return true
}
