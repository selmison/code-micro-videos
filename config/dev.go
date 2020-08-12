// +build dev

package config

import "fmt"

func GetConfig() (*Config, error) {
	dbConnStr = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPass,
		dbSSLMode,
	)
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
