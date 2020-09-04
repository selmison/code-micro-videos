package config

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/testcontainers/testcontainers-go"

	"github.com/selmison/code-micro-videos/pkg/storage/files"
)

const (
	addressServer  = "127.0.0.1:3333"
	containerImage = "postgres:12.3-alpine"
	dbDrive        = "postgres"
	dbName         = "code-micro-videos"
	dbHost         = "127.0.0.1"
	dbPort         = 5432
	dbUser         = "postgres"
	dbPass         = "postgres"
	dbSSLMode      = "disable"
)

var (
	ProjectPath string
)

type Config struct {
	ctx           context.Context
	container     *testcontainers.Container
	AddressServer string
	DBDrive       string
	DBName        string
	DBPort        int
	DBUser        string
	DBPass        string
	DBSSLMode     string
	DBConnStr     string
	RepoFiles     files.Repository
}

func init() {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatalln(fmt.Sprintf(`Error on getting the base path: %s - %s`, err.Error(), string(cmdOut)))
	}
	ProjectPath = strings.TrimSpace(string(cmdOut))
}
