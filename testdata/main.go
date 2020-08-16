package testdata

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
)

var (
	VideoFactory = factory.NewFactory(
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
			default:
				return int16(crud.FreeRating)
			}
		}(), nil
	}).Attr("Duration", func(args factory.Args) (interface{}, error) {
		return int16(randomdata.Number(1, 300)), nil
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
	FakeVideos         []models.Video
	FakeCategoriesDTO  []crud.CategoryDTO
	FakeGenresDTO      []crud.GenreDTO
	FakeCastMembersDTO []crud.CastMemberDTO
	FakeVideosDTO      []crud.VideoDTO
	FakeVideoSlice     models.VideoSlice
)

func init() {
	length := 10
	FakeVideos = make([]models.Video, length)
	for i := 0; i < length; i++ {
		FakeVideos[i] = *(VideoFactory.MustCreate().(*models.Video))
	}
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
	FakeVideosDTO = make([]crud.VideoDTO, len(FakeVideos))
	for i, video := range FakeVideos {
		FakeVideosDTO[i] = crud.VideoDTO{
			Title:        video.Title,
			Description:  video.Description,
			YearLaunched: video.YearLaunched,
			Opened:       video.Opened.Bool,
			Rating:       crud.VideoRating(video.Rating),
			Duration:     video.Duration,
		}
	}
	FakeVideoSlice = make([]*models.Video, len(FakeVideos))
	for i, video := range FakeVideos {
		FakeVideoSlice[i] = &models.Video{
			Title:        video.Title,
			Description:  video.Description,
			YearLaunched: video.YearLaunched,
			Opened:       video.Opened,
			Rating:       video.Rating,
			Duration:     video.Duration,
		}
	}

}
