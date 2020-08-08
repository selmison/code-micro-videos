package seeds

//func InitDB() (*sql.DB, error) {
//	migrations := &migrate.FileMigrationSource{
//		Dir: config.ProjectPath + string(os.PathSeparator) + "migrations",
//	}
//	const (
//		drive    = "sqlite3"
//		filename = "file::memory:?cache=shared"
//	)
//	db, err := sql.Open(drive, filename)
//	if err != nil {
//		return nil, err
//	}
//	n, err := migrate.Exec(db, drive, migrations, migrate.Up)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Printf("Applied %d migrations!\n", n)
//	return db, nil
//}

//func InitDB(db *sql.DB) error {
//	amount := 10
//	ss := MakeSeeds(amount)
//	up := make([]string, amount)
//	down := make([]string, amount)
//	for i, s := range ss {
//		var insert string
//		uuid := uuid2.New()
//		if len(s.Description) == 0 {
//			insert = fmt.Sprintf("INSERT INTO categories (id, name) VALUES (%q, %q);", uuid, s.Name)
//		} else {
//			insert = fmt.Sprintf("INSERT INTO categories (id, name, description) VALUES (%q, %q, %q);", uuid, s.Name, s.Description)
//		}
//		del := fmt.Sprintf("DELETE FROM categories WHERE name = %q;", s.Name)
//		up[i] = insert
//		down[i] = del
//	}
//	fileMigrations := &migrate.FileMigrationSource{
//		Dir: config.ProjectPath + string(os.PathSeparator) + "migrations",
//	}
//	fileMigrationArray, err := fileMigrations.FindMigrations()
//	if err != nil {
//		return err
//	}
//	var migrationArray []*migrate.Migration
//	migrationArray = append(migrationArray, fileMigrationArray...)
//	migrationArray = append(migrationArray, &migrate.Migration{
//		Id:   "20200628135781",
//		Up:   up,
//		Down: down,
//	})
//	migrations := &migrate.MemoryMigrationSource{
//		Migrations: migrationArray,
//	}
//	n, err := migrate.Exec(db, config.DBDrive, migrations, migrate.Up)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("Applied %d migrations!\n", n)
//	return nil
//}

//func ClearCategoriesTable(db *sql.DB) error {
//	_, err := models.Categories().DeleteAll(context.Background(), db)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func ClearCategoriesTable(dbConnStr string) error {
//	db, err := sql.Open(config.DBDrive, dbConnStr)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			log.Fatalln(err)
//		}
//	}()
//	if _, err := models.Categories().DeleteAll(context.Background(), db); err != nil {
//		return err
//	}
//	return nil
//}
