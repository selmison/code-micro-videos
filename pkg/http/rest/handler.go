package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/selmison/code-micro-videos/pkg/listing"
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func Handler(address string, l listing.Service, m modifying.Service) {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})
	CategoryHandler(r, l, m)
	fmt.Printf("The server is on tap now: http://%s\n", address)
	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Fatalln(err)
	}
}
