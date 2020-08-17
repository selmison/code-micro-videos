package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
)

type server struct {
	router *httprouter.Router
	svc    crud.Service
	logger *zap.SugaredLogger
}

func InitApp(ctx context.Context, dbConnStr string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("test: failed to get config: %v\n", err)
	}
	db, err := sql.Open(cfg.DBDrive, dbConnStr)
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
	return initHttpServer(cfg.AddressServer, svc)
}

func initHttpServer(address string, crud crud.Service) error {
	s := newServer(crud)
	fmt.Printf("The server is on tap now: http://%s\n", address)
	if err := http.ListenAndServe(address, s); err != nil {
		return err
	}
	return nil
}

func initLogger() *zap.SugaredLogger {
	devConfig := zap.NewDevelopmentConfig()
	devConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger, err := devConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := zapLogger.Sugar()
	defer func() {
		_ = sugar.Sync()
	}()
	return sugar
}

func newServer(svc crud.Service) *server {
	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := fmt.Fprint(w, "Welcome!\n"); err != nil {
			log.Println(err)
		}
	})
	s := &server{router: r, svc: svc, logger: initLogger()}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Info(r.Method, r.URL.Path)
	s.router.ServeHTTP(w, r)
}

func (s *server) bodyToStruct(w http.ResponseWriter, r *http.Request, dto interface{}) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.errInternalServer(w, err)
		return err
	}
	if err := r.Body.Close(); err != nil {
		s.errInternalServer(w, err)
		return err
	}
	if err := json.Unmarshal(bytes, &dto); err != nil {
		s.errUnprocessableEntity(w, err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			s.errInternalServer(w, err)
			return err
		}
		return err
	}
	return nil
}

func (s *server) errBadRequest(w http.ResponseWriter, err error) {
	s.logger.Info(err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (s *server) errInternalServer(w http.ResponseWriter, err error) {
	s.logger.Error(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *server) errNotFound(w http.ResponseWriter, err error) {
	s.logger.Info(err)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *server) errUnprocessableEntity(w http.ResponseWriter, err error) {
	s.logger.Info(err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

func (s *server) errStatusConflict(w http.ResponseWriter, err error) {
	s.logger.Info(err)
	http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
}
