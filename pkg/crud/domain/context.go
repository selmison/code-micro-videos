package domain

import (
	"context"

	"github.com/selmison/code-micro-videos/pkg/common/domain"
)

type Context interface {
	context.Context
	Logger() domain.Logger
	Repo() Repository
}
