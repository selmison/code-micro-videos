//go:generate sqlboiler psql --no-tests --add-global-variants --add-soft-deletes --config /Users/selmison/Projects/code-micro-videos/backend/sqlboiler.toml --output /Users/selmison/Projects/code-micro-videos/backend/models

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := cfg.TerminateContainer(); err != nil {
			log.Println(err)
		}
	}()
	svc, err := initService(context.Background(), &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return categories.HandlerFromMux(svc, router)
	})
}

func initService(ctx context.Context, cfg *config.Config) (service.Service, error) {
	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	r := sqlboiler.NewRepository(ctx, db, cfg.RepoFiles)
	return service.NewService(r), nil
}
