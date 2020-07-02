package seeds

import (
	"context"
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/selmison/code-micro-videos/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/boil"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Seed struct {
	Name        string
	Description string
}

func init() {
	faker.SetGenerateUniqueValues(true)
}

func createCategory(ctx context.Context, name, description string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	c := models.Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: null.String{String: description, Valid: true},
		CreatedAt:   null.String{String: time.Now().String(), Valid: true},
		UpdatedAt:   null.String{String: time.Now().String(), Valid: true},
	}
	err := c.InsertG(ctx, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func MakeSeeds(amount int) []Seed {
	seeds := make([]Seed, amount)
	for i := 0; i < amount; i++ {
		seeds[i] = Seed{
			Name:        strings.ToLower(faker.FirstName()),
			Description: Sentence(),
		}
	}
	return seeds
}

func Run(ctx context.Context, db *sql.DB, seeds []Seed) {
	boil.SetDB(db)
	for _, seed := range seeds {
		if err := createCategory(ctx, seed.Name, seed.Description); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}

func Sentence() string {
	if rand.Intn(2) == 0 {
		return faker.Sentence()
	}
	return ""
}
