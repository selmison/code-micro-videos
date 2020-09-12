package domain

type CastMemberRepository interface {
	CreateCastMember(ctx Context, castMember CastMember) error
	FetchCastMember(ctx Context, name string) (CastMember, error)
	GetCastMembers(ctx Context, limit int) ([]CastMember, error)
	RemoveCastMember(ctx Context, name string) error
	UpdateCastMember(ctx Context, name string, castMember CastMember) error
}
