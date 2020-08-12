package seeds

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bxcodec/faker/v3"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/selmison/code-micro-videos/config"
)

type Seed struct {
	Name        string
	Description string
}

//func init() {
//	faker.SetGenerateUniqueValues(true)
//}

//func createCategory(ctx context.Context, name, description string) error {
//	name = strings.TrimSpace(strings.ToLower(name))
//	c := models.Category{
//		ID:          uuid.New().String(),
//		Name:        name,
//		Description: null.String{String: description, Valid: true},
//		CreatedAt:   null.Time{Time: time.Now(), Valid: true},
//		UpdatedAt:   null.Time{Time: time.Now(), Valid: true},
//	}
//	err := c.InsertG(ctx, boil.Infer())
//	if err != nil {
//		return err
//	}
//	return nil
//}

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

//func Run(ctx context.Context, db *sql.DB, seeds []Seed) {
//	boil.SetDB(db)
//	for _, seed := range seeds {
//		if err := createCategory(ctx, seed.Name, seed.Description); err != nil {
//			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
//		}
//	}
//}

func Sentence() string {
	if rand.Intn(2) == 0 {
		return faker.Sentence()
	}
	return ""
}

func ApplyMigrations(dbDriver, dbConnStr string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: config.ProjectPath + string(os.PathSeparator) + "migrations",
	}
	db, err := sql.Open(dbDriver, dbConnStr)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	for {
		err = db.Ping()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	n, err := migrate.Exec(db, dbDriver, migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}
