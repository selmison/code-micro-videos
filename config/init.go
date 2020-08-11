package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbDrive   = "postgres"
	dbName    = "code-micro-videos"
	dbPort    = 5432
	dbUser    = "postgres"
	dbPass    = "postgres"
	dbSSLMode = "disable"
)

type DBContainer struct {
	ConnStr string
}

func InitDBContainer(ctx context.Context) (*testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:12.3-alpine",
		Tmpfs: map[string]string{
			"/var/lib/postgresql/data": "rw",
		},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPass,
			"POSTGRES_DB":       dbName,
		},
		ExposedPorts: []string{strconv.Itoa(dbPort)},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}
	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("init dbContainer: %s\n", err)
	}
	return &dbContainer, nil
}

func InitDB(dbConnStr string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: ProjectPath + string(os.PathSeparator) + "migrations",
	}
	db, err := sql.Open(dbDrive, dbConnStr)
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
	n, err := migrate.Exec(db, dbDrive, migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

//func NewCFG() (Config, error) {
//	if cfg == nil {
//		dbContainer, err := InitDBContainer(ctx)
//		if err != nil {
//			return Config{}, err
//		}
//		host, err := (*dbContainer).Host(ctx)
//		if err != nil {
//			return Config{}, fmt.Errorf("access dbContainer: %s\n", err)
//		}
//		port, err := nat.NewPort("tcp", strconv.Itoa(dbPort))
//		if err != nil {
//			return Config{}, err
//		}
//		mappedPort, err := (*dbContainer).MappedPort(ctx, port)
//		if err != nil {
//			return Config{}, fmt.Errorf("access dbContainer: %s\n", err)
//		}
//		dbConnStr := fmt.Sprintf(
//			"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
//			host,
//			mappedPort.Int(),
//			dbName,
//			dbUser,
//			dbPass,
//			dbSSLMode,
//		)
//		cfg = &Config{
//			dbContainer,
//			"127.0.0.1:3333",
//			"postgres",
//			"code-micro-videos",
//			5432,
//			"postgres",
//			"postgres",
//			"disable",
//			dbConnStr,
//		}
//	}
//	return *cfg, nil
//}
