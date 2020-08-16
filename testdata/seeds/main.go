package seeds

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	migrate "github.com/rubenv/sql-migrate"

	"github.com/selmison/code-micro-videos/config"
)

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
