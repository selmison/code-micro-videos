// +build dev

package seeds

import (
	"database/sql"
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/rubenv/sql-migrate"
	"github.com/selmison/code-micro-videos/config"
	"log"
)

func InitDB() {
	amount := 10
	ss := MakeSeeds(amount)
	up := make([]string, amount)
	down := make([]string, amount)
	for i, s := range ss {
		var insert string
		uuid := uuid2.New()
		if len(s.Description) == 0 {
			insert = fmt.Sprintf("INSERT INTO categories (id, name) VALUES (%q, %q);", uuid, s.Name)
		} else {
			insert = fmt.Sprintf("INSERT INTO categories (id, name, description) VALUES (%q, %q, %q);", uuid, s.Name, s.Description)
		}
		delete := fmt.Sprintf("DELETE FROM categories WHERE name = %q;", s.Name)
		up[i] = insert
		down[i] = delete
	}
	fileMigrations := &migrate.FileMigrationSource{
		Dir: "/Users/selmison/Projects/code-micro-videos/backend/migrations/",
	}
	fileMigrationArray, err := fileMigrations.FindMigrations()
	if err != nil {
		log.Fatalln(err)
	}
	var migrationArray []*migrate.Migration
	migrationArray = append(migrationArray, fileMigrationArray...)
	migrationArray = append(migrationArray, &migrate.Migration{
		Id:   "20200628135781",
		Up:   up,
		Down: down,
	})
	migrations := &migrate.MemoryMigrationSource{
		Migrations: migrationArray,
	}
	db, err := sql.Open(config.Drive, config.Url)
	if err != nil {
	}
	n, err := migrate.Exec(db, config.Drive, migrations, migrate.Up)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
