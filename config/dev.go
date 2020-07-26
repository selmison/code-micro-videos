// +build dev

package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	AddressServer = "127.0.0.1:3333"
	Drive         = "sqlite3"
	file          = ".dbdata/db.sqlite"
)

var (
	Url         string
	ProjectPath string
)

func init() {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatalln(fmt.Sprintf(`Error on getting the base path: %s - %s`, err.Error(), string(cmdOut)))
	}
	ProjectPath = strings.TrimSpace(string(cmdOut))
	Url = ProjectPath + string(os.PathSeparator) + file
}
