//+build integration

package config

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

var once sync.Once
var singleInstance *Config

type Config struct {
	ctx           context.Context
	Container     *testcontainers.Container
	AddressServer string
	DBDrive       string
	DBName        string
	DBPort        int
	DBUser        string
	DBPass        string
	DBSSLMode     string
	DBConnStr     string
}

func NewCFG() (*Config, error) {
	var e error
	if singleInstance == nil {
		once.Do(
			func() {
				ctx := context.Background()
				dbContainer, err := InitDBContainer(ctx)
				if err != nil {
					e = err
					return
				}
				host, err := (*dbContainer).Host(ctx)
				if err != nil {
					e = fmt.Errorf("access dbContainer: %s\n", err)
					return
				}
				port, err := nat.NewPort("tcp", strconv.Itoa(dbPort))
				if err != nil {
					e = err
					return
				}
				mappedPort, err := (*dbContainer).MappedPort(ctx, port)
				if err != nil {
					e = fmt.Errorf("access dbContainer: %s\n", err)
					return
				}
				dbConnStr := fmt.Sprintf(
					"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
					host,
					mappedPort.Int(),
					dbName,
					dbUser,
					dbPass,
					dbSSLMode,
				)
				singleInstance = &Config{
					ctx,
					dbContainer,
					"127.0.0.1:3333",
					"postgres",
					"code-micro-videos",
					5432,
					"postgres",
					"postgres",
					"disable",
					dbConnStr,
				}
			})
	}
	if e != nil {
		return nil, e
	}
	return singleInstance, nil
}

func (c *Config) ContainerTerminate() error {
	//cfg = nil
	//if err := (*c.Container).Terminate(c.ctx); err != nil {
	//	return fmt.Errorf("terminate dbContainer: %s", err)
	//}
	return nil
}
