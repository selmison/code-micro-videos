// +build integration

package rest_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/api/rest"
	"github.com/selmison/code-micro-videos/testdata"
	"github.com/selmison/code-micro-videos/testdata/seeds"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	teardownTestCase, err := setupTestMain()
	if err != nil {
		return 1
	}
	defer teardownTestCase(m)
	cfg, err := config.GetConfig()
	if err != nil {
		return 1
	}
	if err := seeds.ApplyMigrations(cfg.DBDrive, cfg.DBConnStr); err != nil {
		log.Fatalln(err, "init db")
		return 1
	}
	var code int
	go func() {
		if err := rest.InitApp(context.Background(), cfg.DBConnStr); err != nil {
			code = 1
			log.Fatalln(err, "init app")
		}
	}()
	if code > 0 {
		return code
	}
	time.Sleep(1 * time.Second)
	code = m.Run()
	return code
}

func setupTestMain() (func(m *testing.M), error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("test: failed to get config: %v", err)
	}
	return func(m *testing.M) {
		if err := testdata.ClearTables(cfg.DBDrive, cfg.DBConnStr); err != nil {
			log.Printf("test: clear categories table: %v/n", err)
		}
	}, nil
}

func setupTestCase(t *testing.T, fakes interface{}) (*config.Config, func(t *testing.T), error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("test: failed to get config: %v", err)
	}
	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
	if err != nil {
		return nil, nil, fmt.Errorf("test: failed to open DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v", err)
		}
	}()
	ctx := context.Background()
	switch v := fakes.(type) {
	case []models.Category:
		for _, category := range v {
			err = category.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, fmt.Errorf("test: insert category: %s", err)
			}
		}
	case []models.Genre:
		for _, genre := range v {
			err = genre.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, fmt.Errorf("test: insert genre: %s", err)
			}
		}
	case []models.CastMember:
		for _, castMember := range v {
			err = castMember.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, fmt.Errorf("test: insert cast member: %s", err)
			}
		}
	case []models.Video:
		for _, video := range v {
			err = video.InsertG(ctx, boil.Infer())
			if err != nil {
				return nil, nil, fmt.Errorf("test: insert video: %s", err)
			}
			err = video.SetCategoriesG(ctx, true, video.R.Categories...)
			if err != nil {
				return nil, nil, fmt.Errorf(
					"test: Insert new a group of categories and assign them to the video: %s",
					err,
				)
			}
			err = video.SetGenresG(ctx, true, video.R.Genres...)
			if err != nil {
				return nil, nil, fmt.Errorf(
					"test: Insert new a group of genres and assign them to the video: %s",
					err,
				)
			}
		}
	}

	return cfg, func(t *testing.T) {
		if err := testdata.ClearTables(cfg.DBDrive, cfg.DBConnStr); err != nil {
			t.Errorf("test: clear categories table: %v", err)
		}
	}, nil
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func toJSON(i interface{}) []byte {
	s, _ := json.Marshal(i)
	return s
}
