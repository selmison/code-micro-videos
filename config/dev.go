// +build dev

package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/selmison/code-micro-videos/models"
)

const (
	AddressServer = "127.0.0.1:3333"
	DBDrive       = "postgres"
	DBName        = "code-micro-videos"
	DBHost        = "localhost"
	DBPort        = 5432
	DBUser        = "postgres"
	DBPass        = "postgres"
	DBSSLMode     = "disable"
)

var (
	DBConnStr = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		DBHost,
		DBPort,
		DBName,
		DBUser,
		DBPass,
		DBSSLMode,
	)
)

func ClearCategoriesTable(dbDriver, dbConnStr string) error {
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
	return nil
}
