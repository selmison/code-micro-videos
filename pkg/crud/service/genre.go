package service

//type Genre struct {
//	Id         *string
//	Name       string
//	Categories *[]CategoryOfGenre
//}
//
//type CategoryOfGenre struct {
//	Id          *string
//	Name        string
//	Description *string
//}

//func (g Genre) RemoveGenreByName(ctx domain.Context, name string) error {
//	name = strings.ToLower(strings.TrimSpace(name))
//	if len(name) == 0 {
//		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
//	}
//	if err := ctx.Repo().RemoveGenreByName(ctx, name); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
//		}
//		return err
//	}
//	return nil
//}
//
//func (g Genre) UpdateGenre(ctx domain.Context, name string, genre Genre) error {
//	name = strings.ToLower(strings.TrimSpace(name))
//	if len(name) == 0 {
//		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
//	}
//	if genre == (Genre{}) {
//		return fmt.Errorf("genres %w", logger.ErrIsEmpty)
//	}
//	validatedName, err := domain.ValidateAndParseValidatableOfGenre(genre.Name)
//	if err != nil {
//		return err
//	}
//	genre.Name = validatedName
//	if err := ctx.Repo().UpdateGenre(ctx, name, genre); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (g Genre) CreateGenre(ctx domain.Context, genre Genre) error {
//	if genre == (Genre{}) {
//		return fmt.Errorf("genres %w", logger.ErrIsEmpty)
//	}
//	validatedName, err := domain.ValidateAndParseValidatableOfGenre(genre.Name)
//	if err != nil {
//		return err
//	}
//	genre.Name = validatedName
//	return ctx.Repo().CreateGenre(ctx, genre)
//}
//
//func (g Genre) GetGenres(ctx domain.Context, limit int) ([]Genre, error) {
//	if limit < 0 {
//		return nil, logger.ErrInvalidedLimit
//	}
//	return ctx.Repo().GetGenres(ctx, limit)
//}
//
//func (g Genre) FetchGenre(ctx domain.Context, name string) (Genre, error) {
//	name = strings.ToLower(strings.TrimSpace(name))
//	genre, err := ctx.Repo().FetchGenre(ctx, name)
//	if err == sql.ErrNoRows {
//		return Genre{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
//	} else if err != nil {
//		return Genre{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
//	}
//	return genre, nil
//}
