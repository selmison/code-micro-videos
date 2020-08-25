package testdata

import (
	"database/sql"
	"log"
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
	if _, err = db.Exec("DELETE FROM videos"); err != nil {
		return err
	}
	if _, err = db.Exec("DELETE FROM categories"); err != nil {
		return err
	}
	if _, err = db.Exec("DELETE FROM genres"); err != nil {
		return err
	}
	if _, err = db.Exec("DELETE FROM cast_members"); err != nil {
		return err
	}
	return nil
}
