package domain

type CategoryRepository interface {
	CreateCategory(ctx Context, category Category) error
	FetchCategory(ctx Context, name string) (Category, error)
	GetCategories(ctx Context, limit int) ([]Category, error)
	RemoveCategory(ctx Context, name string) error
	UpdateCategory(ctx Context, name string, category Category) error
}
