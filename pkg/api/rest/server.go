package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
)

type server struct {
	router *httprouter.Router
	svc    crud.Service
	log    zerolog.Logger
}

func InitApp(ctx context.Context, dbConnStr string) error {
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	r := sqlboiler.NewRepository(ctx, db)
	svc := crud.NewService(r)
	return initHttpServer(config.AddressServer, svc)
}

func initHttpServer(address string, crud crud.Service) error {
	s := newServer(crud)
	fmt.Printf("The server is on tap now: http://%s\n", address)
	if err := http.ListenAndServe(address, s); err != nil {
		return err
	}
	return nil
}

func newServer(svc crud.Service) *server {
	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := fmt.Fprint(w, "Welcome!\n"); err != nil {
			log.Println(err)
		}
	})
	s := &server{router: r, svc: svc}
	s.routes()
	s.logger()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) bodyToStruct(w http.ResponseWriter, r *http.Request, dto interface{}) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.errInternalServer(w, err)
	}
	if err := r.Body.Close(); err != nil {
		s.errInternalServer(w, err)
	}
	if err := json.Unmarshal(bytes, &dto); err != nil {
		s.errUnprocessableEntity(w, err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			s.errInternalServer(w, err)
		}
	}
}

func (s *server) logger() http.Handler {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	c := alice.New()
	c = c.Append(hlog.NewHandler(logger))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	return c.Then(s)
}

func (s *server) errBadRequest(w http.ResponseWriter, err error) {
	s.log.Error().Err(err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (s *server) errInternalServer(w http.ResponseWriter, err error) {
	s.log.Error().Err(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *server) errNotFound(w http.ResponseWriter, err error) {
	s.log.Error().Err(err)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *server) errUnprocessableEntity(w http.ResponseWriter, err error) {
	s.log.Error().Err(err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

func (s *server) errStatusConflict(w http.ResponseWriter, err error) {
	s.log.Error().Err(err)
	http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
}
