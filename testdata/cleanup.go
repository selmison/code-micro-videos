package testdata

import (
	"context"
	"database/sql"
	"log"

	"github.com/selmison/code-micro-videos/models"
)

func ClearTables(dbDriver, dbConnStr string) error {
	db, err := sql.Open(dbDriver, dbConnStr)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	if _, err := models.Categories().DeleteAll(context.Background(), db); err != nil {
		return err
	}
	if _, err := models.Genres().DeleteAll(context.Background(), db, true); err != nil {
		return err
	}
	return nil
}
