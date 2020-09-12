package domain

type GenreRepository interface {
	CreateGenre(ctx Context, genre Genre) error
	FetchGenre(ctx Context, name string) (Genre, error)
	GetGenres(ctx Context, limit int) ([]Genre, error)
	RemoveGenre(ctx Context, name string) error
	UpdateGenre(ctx Context, name string, genre Genre) error
}
