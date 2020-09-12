package domain

type Repository interface {
	CastMemberRepository
	CategoryRepository
	GenreRepository
	VideoRepository
}
