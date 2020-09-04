//go:generate sqlboiler psql --no-tests --add-global-variants --add-soft-deletes --config /Users/selmison/Projects/code-micro-videos/backend/sqlboiler.toml --output /Users/selmison/Projects/code-micro-videos/backend/models

package main

import (
	"context"
	"log"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/pkg/api/rest"
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
	if err := rest.InitApp(context.Background(), &cfg); err != nil {
		log.Fatalln(err)
	}
}
