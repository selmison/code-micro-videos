package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker/v3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"github.com/selmison/code-micro-videos/testdata/seeds"
	"github.com/volatiletech/null/v8"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

var (
	ctx        = context.Background()
	seedsArray []seeds.Seed
)

func init() {
	faker.SetGenerateUniqueValues(true)
	seedsArray = seeds.MakeSeeds(10)
}

func TestRepository_AddCategory(t *testing.T) {
	fakeName := faker.FirstName()
	fakeDescription := seeds.Sentence()
	type args struct {
		c modifying.CategoryDTO
	}
	type returns struct {
		c models.Category
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name in CategoryDTO already exists",
			args: args{modifying.CategoryDTO{
				Name: seedsArray[0].Name,
			}},
			want:    returns{models.Category{}, modifying.ErrAlreadyExists},
			wantErr: true,
		},
		{
			name: "When CategoryDTO is right",
			args: args{
				modifying.CategoryDTO{
					Name:        fakeName,
					Description: fakeDescription,
				},
			},
			want: returns{
				models.Category{
					Name:        fakeName,
					Description: null.String{String: fakeDescription, Valid: true},
				},
				nil,
			},
			wantErr: false,
		},
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("test: failed to open DB: %v\n", err)
	}
	defer db.Close()
	seeds.Run(ctx, db, seedsArray)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(ctx, db)
			err := r.AddCategory(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_GetCategories(t *testing.T) {
	type args struct {
		limit int
	}
	fakeLimit := 10
	type returns struct {
		c models.CategorySlice
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When limit is right",
			args:    args{fakeLimit},
			want:    returns{nil, nil},
			wantErr: false,
		},
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("test: failed to open DB: %v\n", err)
	}
	defer db.Close()
	seeds.Run(ctx, db, seedsArray)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(ctx, db)
			got, err := r.GetCategories(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != fakeLimit {
				t.Errorf("GetCategories() len(got): %v, want: %d", got, fakeLimit)
			}
		})
	}
}

func TestRepository_FetchCategory(t *testing.T) {
	indexRandom := rand.Intn(10)
	fakeNames := [2]string{
		faker.FirstName(),
		seedsArray[indexRandom].Name,
	}
	type args struct {
		name string
	}
	type returns struct {
		c models.Category
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name is not found",
			args: args{fakeNames[0]},
			want: returns{
				models.Category{},
				sql.ErrNoRows,
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeNames[1]},
			want: returns{
				models.Category{
					Name: fakeNames[1],
				},
				nil,
			},
			wantErr: false,
		},
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("test: failed to open DB: %v\n", err)
	}
	defer db.Close()
	seeds.Run(ctx, db, seedsArray)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(ctx, db)
			got, err := r.FetchCategory(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.c.Name {
				t.Errorf("FetchCategory() got: %q, want: %q", got.Name, tt.want.c.Name)
			}
		})
	}
}

func TestRepository_RemoveCategory(t *testing.T) {
	indexRandom := rand.Intn(10)
	fakeNames := [2]string{
		faker.FirstName(),
		seedsArray[indexRandom].Name,
	}
	type fields struct {
		ctx context.Context
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "When name is not found",
			args:    args{fakeNames[0]},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name:    "When name is found",
			args:    args{fakeNames[1]},
			want:    nil,
			wantErr: false,
		},
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("test: failed to open DB: %v\n", err)
	}
	defer db.Close()
	seeds.Run(ctx, db, seedsArray)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(ctx, db)
			err := r.RemoveCategory(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveCategory() got: %s, want: %q", err, tt.want)
			}
		})
	}
}

func TestRepository_UpdateCategory(t *testing.T) {
	indexRandom := rand.Intn(len(seedsArray) - 1)
	fakeNames := map[string]string{
		"inedited": faker.UUIDDigit(),
		"seed1":    seedsArray[indexRandom].Name,
		"seed2":    seedsArray[indexRandom+1].Name,
	}
	type fields struct {
		ctx context.Context
	}
	type args struct {
		name string
		c    modifying.CategoryDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When name in CategoryDTO already exists",
			args: args{
				fakeNames["seed1"],
				modifying.CategoryDTO{
					Name:        fakeNames["seed2"],
					Description: faker.Sentence(),
				},
			},
			want:    fmt.Errorf("%s %w", strings.TrimSpace(fakeNames["seed2"]), modifying.ErrAlreadyExists),
			wantErr: true,
		},
		{
			name: "When name exists and CategoryDTO is right",
			args: args{
				fakeNames["seed1"],
				modifying.CategoryDTO{
					Name:        fakeNames["inedited"],
					Description: faker.Sentence(),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("test: failed to open DB: %v\n", err)
	}
	defer db.Close()
	seeds.Run(ctx, db, seedsArray)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(ctx, db)
			err := r.UpdateCategory(tt.args.name, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateCategory() got: %v, want: %v", err, tt.want)
			}
		})
	}
}

func initDB() (*sql.DB, error) {
	migrations := &migrate.FileMigrationSource{
		Dir: "/Users/selmison/Projects/code-micro-videos/backend/migrations/",
	}
	const (
		drive    = "sqlite3"
		filename = "file::memory:?cache=shared"
	)
	db, err := sql.Open(drive, filename)
	if err != nil {
		return nil, err
	}
	n, err := migrate.Exec(db, drive, migrations, migrate.Up)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return db, nil
}
