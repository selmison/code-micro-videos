package testdata

import (
	"log"
	"math/rand"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"github.com/google/uuid"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/category"
	"github.com/selmison/code-micro-videos/pkg/genre"
	"github.com/selmison/code-micro-videos/pkg/video"
)

const (
	FakeCastMembersLength = 10
	FakeCategoriesLength  = 10
	FakeGenresLength      = 10
	FakeVideosLength      = 10
)

var (
	newCastMemberFactory = factory.NewFactory(
		&cast_member.NewCastMemberDTO{},
	).Attr("Name", func(args factory.Args) (interface{}, error) {
		return strings.TrimSpace(randomdata.FullName(randomdata.RandomGender)), nil
	}).Attr("Type", func(args factory.Args) (interface{}, error) {
		return func() cast_member.CastMemberType {
			switch randomdata.Number(2) {
			case 0:
				return cast_member.Actor
			}
			return cast_member.Director
		}(), nil
	})

	categoryFactory = factory.NewFactory(
		&category.Category{},
	).Attr("Id", func(args factory.Args) (interface{}, error) {
		return uuid.New().String(), nil
	}).Attr("Name", func(args factory.Args) (interface{}, error) {
		return strings.ToLower(randomdata.SillyName()), nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		desc := ""
		if rand.Intn(2) == 0 {
			desc = randomdata.Paragraph()
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
		&genre.Genre{},
	).Attr("Id", func(args factory.Args) (interface{}, error) {
		return uuid.New().String(), nil
	}).Attr("Name", func(args factory.Args) (interface{}, error) {
		return strings.ToLower(randomdata.SillyName()), nil
	}).Attr("IsValidated", func(args factory.Args) (interface{}, error) {
		isValidated := true
		if rand.Intn(2) == 0 {
			isValidated = false
		}
		return isValidated, nil
	})

	videoFactory = factory.NewFactory(
		&video.Video{},
	).Attr("Id", func(args factory.Args) (interface{}, error) {
		return uuid.New().String(), nil
	}).Attr("Title", func(args factory.Args) (interface{}, error) {
		return strings.TrimSpace(randomdata.FullName(randomdata.RandomGender)), nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		return randomdata.Paragraph(), nil
	}).Attr("YearLaunched", func(args factory.Args) (interface{}, error) {
		return int16(randomdata.Number(1900, 2030)), nil
	}).Attr("Opened", func(args factory.Args) (interface{}, error) {
		return randomdata.Boolean(), nil
	}).Attr("Rating", func(args factory.Args) (interface{}, error) {
		return func() video.VideoRating {
			switch randomdata.Number(6) {
			case 0:
				return video.FreeRating
			case 1:
				return video.TenRating
			case 2:
				return video.TwelveRating
			case 3:
				return video.FourteenRating
			case 4:
				return video.SixteenRating
			case 5:
				return video.EighteenRating
			}
			return video.FreeRating
		}(), nil
	}).Attr("CategoriesId", func(args factory.Args) (interface{}, error) {
		index := randomdata.Number(FakeCategoriesLength)
		return []string{FakeCategories[index].Id}, nil
	}).Attr("GenresId", func(args factory.Args) (interface{}, error) {
		index := randomdata.Number(FakeGenresLength)
		return []string{FakeGenres[index].Id}, nil
	})

	FakeCategories     []category.Category
	FakeNewCategories  []category.NewCategory
	FakeGenres         []genre.Genre
	FakeNewGenres      []genre.NewGenre
	FakeCastMembers    []cast_member.CastMember
	FakeCastMemberDTOs []cast_member.CastMemberDTO
	FakeNewCastMembers []cast_member.NewCastMemberDTO
	FakeVideos         []video.Video
	FakeNewVideos      []video.NewVideo
	FakeVideoSlice     models.VideoSlice
)

func init() {
	var err error
	FakeCastMembers, FakeNewCastMembers, FakeCastMemberDTOs, err = generateFakeCastMembers(FakeCastMembersLength)
	if err != nil {
		log.Fatalf("init: failed to generateFakeCastMembers: %v\n", err)
	}
	FakeCategories, FakeNewCategories = generateFakeCategories(FakeCategoriesLength)
	FakeGenres, FakeNewGenres = generateFakeGenres(FakeGenresLength)
	FakeVideos, FakeNewVideos, FakeVideoSlice = generateFakeVideos(FakeVideosLength)
}

func generateFakeCategories(length int) ([]category.Category, []category.NewCategory) {
	fakeCategories := make([]category.Category, length)
	for i := 0; i < length; i++ {
		fakeCategories[i] = *(categoryFactory.MustCreate().(*category.Category))
	}
	fakeNewCategory := make([]category.NewCategory, length)
	for i, fakeCategory := range fakeCategories {
		fakeNewCategory[i] = category.NewCategory{
			Name:        fakeCategory.Name,
			Description: fakeCategory.Description,
			GenresId:    fakeCategory.GenresId,
			IsValidated: fakeCategory.IsValidated,
		}
	}
	return fakeCategories, fakeNewCategory
}

func generateFakeCastMembers(length int) (
	[]cast_member.CastMember,
	[]cast_member.NewCastMemberDTO,
	[]cast_member.CastMemberDTO,
	error,
) {
	fakeNewCastMembers := make([]cast_member.NewCastMemberDTO, length)
	fakeCastMembers := make([]cast_member.CastMember, length)
	fakeCastMemberDTOs := make([]cast_member.CastMemberDTO, length)
	var err error
	for i := 0; i < length; i++ {
		fakeNewCastMembers[i] = *(newCastMemberFactory.MustCreate().(*cast_member.NewCastMemberDTO))
		fakeCastMembers[i], err = cast_member.NewCastMember(
			uuid.New().String(), fakeNewCastMembers[i])
		if err != nil {
			return nil, fakeNewCastMembers, nil, err
		}
		fakeCastMemberDTOs[i] = cast_member.CastMemberDTO{
			Id:   fakeCastMembers[i].Id(),
			Name: fakeCastMembers[i].Name(),
			Type: fakeCastMembers[i].Type(),
		}
	}
	return fakeCastMembers, fakeNewCastMembers, fakeCastMemberDTOs, nil
}

func generateFakeGenres(length int) ([]genre.Genre, []genre.NewGenre) {
	fakeGenres := make([]genre.Genre, length)
	for i := 0; i < length; i++ {
		fakeGenres[i] = *(genreFactory.MustCreate().(*genre.Genre))
	}
	fakeNewGenre := make([]genre.NewGenre, length)
	for i, fakeGenre := range fakeGenres {
		fakeNewGenre[i] = genre.NewGenre{
			Name:         fakeGenre.Name,
			CategoriesId: fakeGenre.CategoriesId,
			IsValidated:  fakeGenre.IsValidated,
		}
	}
	return fakeGenres, fakeNewGenre
}

func generateFakeVideos(length int) ([]video.Video, []video.NewVideo, models.VideoSlice) {
	fakeVideos := make([]video.Video, length)
	for i := 0; i < length; i++ {
		fakeVideos[i] = *(videoFactory.MustCreate().(*video.Video))
	}
	fakeNewVideos := make([]video.NewVideo, length)
	for i, fakeVideo := range fakeVideos {
		fakeNewVideos[i] = video.NewVideo{
			Title:        fakeVideo.Title,
			Description:  fakeVideo.Description,
			YearLaunched: &fakeVideo.YearLaunched,
			Opened:       fakeVideo.Opened,
			Rating:       &fakeVideo.Rating,
			Duration:     &fakeVideo.Duration,
			CategoriesId: fakeVideo.CategoriesId,
			GenresId:     fakeVideo.GenresId,
		}
	}
	fakeVideoSlice := make([]*models.Video, length)
	//for i, fakeVideo := range fakeVideos {
	//	//yearLaunched := fakeVideo.YearLaunched
	//	//opened := fakeVideo.Opened.Bool
	//	//rating := fakeVideo.VideoRating(fakeVideo.Rating)
	//	//duration := fakeVideo.Duration
	//	fakeVideoSlice[i] = &models.Video{
	//		ID:           fakeVideo.Id,
	//		Title:        fakeVideo.Title,
	//		Description:  fakeVideo.Description,
	//		YearLaunched: *fakeVideo.YearLaunched,
	//		Opened:       fakeVideo.Opened,
	//		Rating:       fakeVideo.Rating,
	//		Duration:     *fakeVideo.Duration,
	//		R:            fakeVideo.R,
	//	}
	//}
	return fakeVideos, fakeNewVideos, fakeVideoSlice
}
