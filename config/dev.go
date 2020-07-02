// +build dev

package config

import (
	"log"
	"os"
)

const (
	AddressServer = "127.0.0.1:3333"
	Drive         = "sqlite3"
	file          = ".dbdata/db.sqlite"
)

var (
	Url string
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	Url = path + string(os.PathSeparator) + file
}
