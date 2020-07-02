package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/pkg/http/rest"
	"github.com/selmison/code-micro-videos/pkg/listing"
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata/seeds"
	"log"
)

func main() {
	seeds.InitDB()
	ctx := context.Background()
	db, err := sql.Open(config.Drive, config.Url)
	if err != nil {
		log.Fatalln(err)
	}

	r := sqlboiler.NewRepository(ctx, db)
	lister := listing.NewService(r)
	modifier := modifying.NewService(r)
	rest.Handler(config.AddressServer, lister, modifier)
}
