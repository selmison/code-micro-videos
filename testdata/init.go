package testdata

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/selmison/code-micro-videos/config"
)

type DBContainer struct {
	ConnStr   string
	Container *testcontainers.Container
}

func NewDBContainer(ctx context.Context) (*DBContainer, error) {
	port, err := nat.NewPort("tcp", strconv.Itoa(config.DBPort))
	if err != nil {
		return nil, err
	}
	req := testcontainers.ContainerRequest{
		Image: "postgres:12.3-alpine",
		Tmpfs: map[string]string{
			"/var/lib/postgresql/data": "rw",
		},
		Env: map[string]string{
			"POSTGRES_USER":     config.DBUser,
			"POSTGRES_PASSWORD": config.DBPass,
			"POSTGRES_DB":       config.DBName,
		},
		ExposedPorts: []string{port.Port()},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}
	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("init dbContainer: %s\n", err)
	}
	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("access dbContainer: %s\n", err)
	}

	mappedPort, err := dbContainer.MappedPort(ctx, port)
	if err != nil {
		return nil, fmt.Errorf("access dbContainer: %s\n", err)
	}
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		host,
		mappedPort.Int(),
		config.DBName,
		config.DBUser,
		config.DBPass,
		config.DBSSLMode,
	)
	return &DBContainer{dbConnStr, &dbContainer}, nil
}
