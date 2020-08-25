package testdata

import (
	"math/rand"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
)

const (
	FakeVideosLength           = 10
	fakeVideosCategoriesLength = 3
	fakeVideosGenresLength     = 3
)

type RemoverOfTimes interface {
	removeTimes()
}

var (
	categoryFactory = factory.NewFactory(
		&models.Category{},
	).Attr("ID", func(args factory.Args) (interface{}, error) {
		return uuid.New().String(), nil
	}).Attr("Name", func(args factory.Args) (interface{}, error) {
		return randomdata.SillyName(), nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		desc := null.String{String: "", Valid: true}
		if rand.Intn(2) == 0 {
			desc = null.String{String: randomdata.Paragraph(), Valid: true}
		}
		return desc, nil
	}).Attr("IsValidated", func(args factory.Args) (interface{}, error) {
		isValidated := true
		if rand.Intn(2) == 0 {
			isValidated = false
		}
		return isValidated, nil
	})

	genreFactory = factory.NewFactory(
		&models.Genre{},
	).Attr("ID", func(args factory.Args) (interface{}, error) {
		return uuid.New().String(), nil
	}).Attr("Name", func(args factory.Args) (interface{}, error) {
		return randomdata.SillyName(), nil
	}).Attr("IsValidated", func(args factory.Args) (interface{}, error) {
		isValidated := true
		if rand.Intn(2) == 0 {
			isValidated = false
		}
		return isValidated, nil
	})

	videoFactory = factory.NewFactory(
		&models.Video{},
	).Attr("ID", func(args factory.Args) (interface{}, error) {
		return uuid.New().String(), nil
	}).Attr("Title", func(args factory.Args) (interface{}, error) {
		return randomdata.FullName(randomdata.RandomGender), nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		return randomdata.Paragraph(), nil
	}).Attr("YearLaunched", func(args factory.Args) (interface{}, error) {
		return int16(randomdata.Number(1900, 2030)), nil
	}).Attr("Opened", func(args factory.Args) (interface{}, error) {
		return null.BoolFrom(randomdata.Boolean()), nil
	}).Attr("Rating", func(args factory.Args) (interface{}, error) {
		return func() int16 {
			switch randomdata.Number(6) {
			case 0:
				return int16(crud.FreeRating)
			case 1:
				return int16(crud.TenRating)
			case 2:
				return int16(crud.TwelveRating)
			case 3:
				return int16(crud.FourteenRating)
			case 4:
				return int16(crud.SixteenRating)
			case 5:
				return int16(crud.EighteenRating)
			}
			return int16(crud.FreeRating)
		}(), nil
	}).Attr("Duration", func(args factory.Args) (interface{}, error) {
		return int16(randomdata.Number(1, 300)), nil
	}).SubFactory("R", videoRFactory)

	videoRFactory = factory.NewFactory(
		models.Video{}.R.NewStruct(),
	).SubSliceFactory("Categories", categoryFactory, func() int {
		return fakeVideosCategoriesLength
	}).SubSliceFactory("Genres", genreFactory, func() int {
		return fakeVideosGenresLength
	})

	FakeCategories = []models.Category{
		{
			ID:          uuid.New().String(),
			Name:        "action",
			Description: null.String{String: "action films", Valid: true},
		},
		{
			ID:   uuid.New().String(),
			Name: "animation",
		},
		{
			ID:          uuid.New().String(),
			Name:        "science fiction",
			Description: null.String{String: "science fiction films", Valid: true},
		},
		{
			ID:   uuid.New().String(),
			Name: "violent",
		},
		{
			ID:          uuid.New().String(),
			Name:        "drama",
			Description: null.String{String: "drama films", Valid: true},
		},
		{
			ID:   uuid.New().String(),
			Name: "romance",
		},
	}
	FakeGenres = []models.Genre{
		{
			ID:   uuid.New().String(),
			Name: "action",
		},
		{
			ID:   uuid.New().String(),
			Name: "animation",
		},
		{
			ID:   uuid.New().String(),
			Name: "science fiction",
		},
		{
			ID:   uuid.New().String(),
			Name: "violent",
		},
		{
			ID:   uuid.New().String(),
			Name: "drama",
		},
		{
			ID:   uuid.New().String(),
			Name: "romance",
		},
	}
	FakeCastMembers = []models.CastMember{
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Director),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Director),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
	}
	FakeCategoriesDTO  []crud.CategoryDTO
	FakeGenresDTO      []crud.GenreDTO
	FakeCastMembersDTO []crud.CastMemberDTO

	FakeVideos     []models.Video
	FakeVideosDTO  []crud.VideoDTO
	FakeVideoSlice models.VideoSlice
)

func init() {
	FakeVideos, FakeVideosDTO, FakeVideoSlice = generateFakeVideos(FakeVideosLength)
	FakeCategoriesDTO = make([]crud.CategoryDTO, len(FakeCategories))
	for i, category := range FakeCategories {
		FakeCategoriesDTO[i] = crud.CategoryDTO{
			Name:        category.Name,
			Description: category.Description.String,
		}
	}
	FakeGenresDTO = make([]crud.GenreDTO, len(FakeGenres))
	for i, user := range FakeGenres {
		FakeGenresDTO[i] = crud.GenreDTO{
			Name: user.Name,
		}
	}
	FakeCastMembersDTO = make([]crud.CastMemberDTO, len(FakeCastMembers))
	for i, castMember := range FakeCastMembers {
		FakeCastMembersDTO[i] = crud.CastMemberDTO{
			Name: castMember.Name,
			Type: crud.CastMemberType(castMember.Type),
		}
	}
}

func generateFakeVideos(length int) ([]models.Video, []crud.VideoDTO, models.VideoSlice) {
	fakeVideos := make([]models.Video, length)
	for i := 0; i < length; i++ {
		fakeVideos[i] = *(videoFactory.MustCreate().(*models.Video))
	}
	fakeVideosDTO := make([]crud.VideoDTO, length)
	for i, video := range fakeVideos {
		categoriesDTO := make([]crud.CategoryDTO, len(video.R.Categories))
		for i, category := range video.R.Categories {
			categoriesDTO[i] = crud.CategoryDTO{
				Name:        category.Name,
				Description: category.Description.String,
			}
		}
		genresDTO := make([]crud.GenreDTO, len(video.R.Genres))
		for i, genre := range video.R.Genres {
			genresDTO[i] = crud.GenreDTO{
				Name: genre.Name,
			}
		}
		yearLaunched := video.YearLaunched
		opened := video.Opened.Bool
		rating := crud.VideoRating(video.Rating)
		duration := video.Duration
		fakeVideosDTO[i] = crud.VideoDTO{
			Title:        video.Title,
			Description:  video.Description,
			YearLaunched: &yearLaunched,
			Opened:       opened,
			Rating:       &rating,
			Duration:     &duration,
			Categories:   categoriesDTO,
			Genres:       genresDTO,
		}
	}
	fakeVideoSlice := make([]*models.Video, length)
	for i, video := range fakeVideos {
		fakeVideoSlice[i] = &models.Video{
			ID:           video.ID,
			Title:        video.Title,
			Description:  video.Description,
			YearLaunched: video.YearLaunched,
			Opened:       video.Opened,
			Rating:       video.Rating,
			Duration:     video.Duration,
			R:            video.R,
		}
	}
	return fakeVideos, fakeVideosDTO, fakeVideoSlice
}
