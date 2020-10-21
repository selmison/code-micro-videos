package inmem

import (
	"context"
	"fmt"
	"testing"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/category"
	"github.com/selmison/code-micro-videos/pkg/genre"
	"github.com/selmison/code-micro-videos/pkg/video"
)

func SetupTestCase(_ *testing.T, fakes interface{}) (func(t *testing.T), cast_member.Repository, error) {
	r := NewCastMemberRepository().(*castMemberRepository)
	switch v := fakes.(type) {
	case []category.Category:
		//for _, category := range v {
		//	err = category.InsertG(ctx, boil.Infer())
		//	if err != nil {
		//		return nil, nil, fmt.Errorf("test: insert category: %s", err)
		//	}
		//}
	case []genre.Genre:
		//for _, genre := range v {
		//	err = genre.InsertG(ctx, boil.Infer())
		//	if err != nil {
		//		return nil, nil, fmt.Errorf("insert genre: %s", err)
		//	}
		//}
	case []cast_member.CastMember:
		for _, castMember := range v {
			if err := r.Store(context.Background(), castMember); err != nil {
				return nil, nil, fmt.Errorf("insert CastMember: %s", err)
			}
			//r.castMembers[castMember.Id()] = castMember
		}
	case []video.Video:
		//for _, video := range v {
		//	err = video.InsertG(ctx, boil.Infer())
		//	if err != nil {
		//		return nil, nil, fmt.Errorf("insert video: %s", err)
		//	}
		//	err = video.SetCategoriesG(ctx, true, video.R.Categories...)
		//	if err != nil {
		//		return nil, nil, fmt.Errorf(
		//			"Insert a new slice of categories and assign them to the video: %s",
		//			err,
		//		)
		//	}
		//	err = video.SetGenresG(ctx, true, video.R.Genres...)
		//	if err != nil {
		//		return nil, nil, fmt.Errorf(
		//			"Insert a new slice of genres and assign them to the video: %s",
		//			err,
		//		)
		//	}
		//}
	}
	return func(t *testing.T) {
		_ = r.DeleteAll(nil)
	}, r, nil
}
