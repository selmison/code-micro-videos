package config

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var (
	ProjectPath string
)

func init() {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatalln(fmt.Sprintf(`Error on getting the base path: %s - %s`, err.Error(), string(cmdOut)))
	}
	ProjectPath = strings.TrimSpace(string(cmdOut))
}
