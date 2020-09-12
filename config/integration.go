//+build integration

package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/selmison/code-micro-videos/pkg/storage/files/memory"
)

func GetConfig() (Config, error) {
	ctx := context.Background()
	//dbContainer, err := InitDBContainer(ctx)
	//if err != nil {
	//	return Config{}, err
	//}
	//host, err := (*dbContainer).Host(ctx)
	//if err != nil {
	//	return Config{}, fmt.Errorf("access dbContainer: %s\n", err)
	//}
	//port, err := nat.NewPort("tcp", strconv.Itoa(dbPort))
	//if err != nil {
	//	return Config{}, err
	//}
	//mappedPort, err := (*dbContainer).MappedPort(ctx, port)
	//if err != nil {
	//	return Config{}, fmt.Errorf("access dbContainer: %s\n", err)
	//}
	//dbConnStr := fmt.Sprintf(
	//	"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
	//	host,
	//	mappedPort.Int(),
	//	dbName,
	//	dbUser,
	//	dbPass,
	//	dbSSLMode,
	//)
	//return Config{
	//	ctx,
	//	dbContainer,
	//	addressServer,
	//	dbDrive,
	//	dbName,
	//	dbPort,
	//	dbUser,
	//	dbPass,
	//	dbSSLMode,
	//	dbConnStr,
	//	memory.NewRepository(),
	//}, nil
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPass,
		dbSSLMode,
	)
	return Config{
		ctx,
		nil,
		addressServer,
		dbDrive,
		dbName,
		dbPort,
		dbUser,
		dbPass,
		dbSSLMode,
		dbConnStr,
		memory.NewRepository(),
	}, nil
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
	//if err := (*c.container).Terminate(c.ctx); err != nil {
	//	return fmt.Errorf("terminate dbContainer: %s", err)
	//}
	return nil
}
