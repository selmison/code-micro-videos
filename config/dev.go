// +build dev

package config

import (
	"fmt"
)

var (
	dbConnStr = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPass,
		dbSSLMode,
	)
)

func NewConfig(addressServer string) *Config {
	return &Config{AddressServer: addressServer}
}

func GetConfig() (*Config, error) {
	return &Config{
		AddressServer: addressServer,
		DBDrive:       dbDrive,
		DBName:        dbName,
		DBPort:        dbPort,
		DBUser:        dbUser,
		DBPass:        dbPass,
		DBSSLMode:     dbSSLMode,
		DBConnStr:     dbConnStr,
	}, nil
}
