//+build integration

package config

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/selmison/code-micro-videos/pkg/storage/files/memory"
)

var (
	once           sync.Once
	singleInstance *Config
)

func GetConfig() (*Config, error) {
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
					addressServer,
					dbDrive,
					dbName,
					dbPort,
					dbUser,
					dbPass,
					dbSSLMode,
					dbConnStr,
					memory.NewRepository(),
				}
			})
	}
	if e != nil {
		return nil, e
	}
	return singleInstance, nil
}

func InitDBContainer(ctx context.Context) (*testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: containerImage,
		Tmpfs: map[string]string{
			"/var/lib/postgresql/data": "rw",
		},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPass,
			"POSTGRES_DB":       dbName,
		},
		ExposedPorts: []string{strconv.Itoa(dbPort)},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}
	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("init dbContainer: %s\n", err)
	}
	return &dbContainer, nil
}

func (c *Config) TerminateContainer() error {
	if err := (*c.container).Terminate(c.ctx); err != nil {
		return fmt.Errorf("terminate dbContainer: %s", err)
	}
	return nil
}
