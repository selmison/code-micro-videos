// +build integration

package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/testdata"
	"github.com/selmison/code-micro-videos/testdata/seeds"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	teardownTestMain, err := setupTestMain()
	if err != nil {
		return 1
	}
	defer teardownTestMain(m)
	cfg, err := config.GetConfig()
	if err != nil {
		return 1
	}
	if err := seeds.ApplyMigrations(cfg.DBDrive, cfg.DBConnStr); err != nil {
		log.Fatalln("init db: ", err)
		return 1
	}
	return m.Run()
}

func setupTestMain() (func(m *testing.M), error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("test: failed to get config: %v\n", err)
	}
	return func(m *testing.M) {
		if err := testdata.ClearTables(cfg.DBDrive, cfg.DBConnStr); err != nil {
			log.Printf("test: clear categories table: %v/n", err)
		}
	}, nil
}

func setupTestCase(fakes interface{}) (*config.Config, func(t *testing.T), *Repository, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("test: failed to get config: %v", err)
	}
	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("test: failed to open DB: %v", err)
	}
	ctx := context.Background()
	r := NewRepository(ctx, db, cfg.RepoFiles)
	switch v := fakes.(type) {
	case []models.Category:
		for _, category := range v {
			err = category.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, nil, fmt.Errorf("test: insert category: %s", err)
			}
		}
	case []models.Genre:
		for _, genre := range v {
			err = genre.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, nil, fmt.Errorf("test: insert genre: %s", err)
			}
		}
	case []models.CastMember:
		for _, castMember := range v {
			err = castMember.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, nil, fmt.Errorf("test: insert cast member: %s", err)
			}
		}
	case []models.Video:
		for _, video := range v {
			err = video.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, nil, fmt.Errorf("test: insert video: %s", err)
			}
			err = video.SetCategoriesG(ctx, true, video.R.Categories...)
			if err != nil {
				return nil, nil, nil, fmt.Errorf(
					"test: Insert a new slice of categories and assign them to the video: %s",
					err,
				)
			}
			err = video.SetGenresG(ctx, true, video.R.Genres...)
			if err != nil {
				return nil, nil, nil, fmt.Errorf(
					"test: Insert a new slice of genres and assign them to the video: %s",
					err,
				)
			}
		}
	}
	return cfg, func(t *testing.T) {
		defer func() {
			if err := db.Close(); err != nil {
				t.Errorf("test: failed to close DB: %v", err)
			}
		}()
		if err := testdata.ClearTables(cfg.DBDrive, cfg.DBConnStr); err != nil {
			t.Errorf("test: clear categories table: %v", err)
		}
	}, r, nil
}
