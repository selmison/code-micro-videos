package genre

import (
	"context"

	"github.com/go-kit/kit/log"
)

// LoggingMiddleware describes a service (as opposed to endpoint) middleware.
type LoggingMiddleware func(Service) Service

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

// NewLoggingMiddleware takes a logger as a dependency
// and returns a service LoggingMiddleware.
func NewLoggingMiddleware(logger log.Logger) LoggingMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

func (m loggingMiddleware) log(method string, input interface{}, output interface{}, err error) {
	if err == nil {
		return
	}
	if input == nil && output == nil {
		_ = m.logger.Log(
			"method:", method,
			"err:", err,
		)
		return
	}
	if input == nil {
		_ = m.logger.Log(
			"method:", method,
			"output:", output,
			"err:", err,
		)
		return
	}
	if output == nil {
		_ = m.logger.Log(
			"method:", method,
			"input:", input,
			"err:", err,
		)
		return
	}
	_ = m.logger.Log(
		"method:", method,
		"input:", input,
		"output:", output,
		"err:", err,
	)
}

func (m loggingMiddleware) Create(ctx context.Context, newGenre NewGenre) (output Genre, err error) {
	defer func() {
		m.log("Create", newGenre, output, err)
	}()
	return m.next.Create(ctx, newGenre)
}

func (m loggingMiddleware) Destroy(ctx context.Context, id string) (err error) {
	defer func() {
		m.log("Destroy", id, nil, err)
	}()
	return m.next.Destroy(ctx, id)
}

func (m loggingMiddleware) List(ctx context.Context) (output []Genre, err error) {
	defer func() {
		m.log("List", nil, output, err)
	}()
	return m.next.List(ctx)
}

func (m loggingMiddleware) Show(ctx context.Context, id string) (output Genre, err error) {
	defer func() {
		m.log("Show", id, output, err)
	}()
	return m.next.Show(ctx, id)
}

func (m loggingMiddleware) Update(ctx context.Context, id string, updateGenre UpdateGenre) (err error) {
	defer func() {
		m.log(
			"Show",
			map[string]interface{}{
				"id":          id,
				"updateGenre": updateGenre,
			},
			nil,
			err,
		)
	}()
	return m.next.Update(ctx, id, updateGenre)
}