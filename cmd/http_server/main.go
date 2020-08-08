//go:generate sqlboiler psql --no-tests --add-global-variants --add-soft-deletes --config /Users/selmison/Projects/code-micro-videos/backend/sqlboiler.toml --output /Users/selmison/Projects/code-micro-videos/backend/models

package main

import (
	"context"
	"log"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/pkg/api/rest"
)

func main() {
	ctx := context.Background()
	log.Fatalln(rest.InitApp(ctx, config.DBConnStr))
}
