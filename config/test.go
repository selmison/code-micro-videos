// +build test

package config

import "os"

const (
	AddressServer = "127.0.0.1:3333"
	Drive         = "sqlite3"
	file          = ".dbdata/db.sqlite"
	DbName        = "code-micro-videos"
	DbHostname    = "db"
	DbPort        = ""
	DbUser        = "postgres"
	DbPass        = "postgres"
)

var (
	DBUrl = ProjectPath + string(os.PathSeparator) + file
)
