package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	kitLog "github.com/go-kit/kit/log"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/category"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/genre"
	"github.com/selmison/code-micro-videos/pkg/id_generator"
	"github.com/selmison/code-micro-videos/pkg/storage/inmem"
	"github.com/selmison/code-micro-videos/pkg/video"
)

type Router interface {
	AddRoute(method string, pattern string, handler http.Handler)
}

type server struct {
	router *httprouter.Router
	svc    service.Service
	//logger domain.Logger
}

func InitHttpServer(address string) error {
	s := NewServer()
	fmt.Printf("The server is on tap now: http://%s\n", address)
	if err := http.ListenAndServe(address, s); err != nil {
		return err
	}
	return nil
}

func NewServer() http.Handler {
	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := fmt.Fprint(w, "Welcome!\n"); err != nil {
			log.Println(err)
		}
	})
	//, logger: domain.NewLogger("debug")
	s := &server{router: r}
	logger := kitLog.NewLogfmtLogger(os.Stderr)
	idGenerator := id_generator.NewGenerator()

	categoryRepo := category.NewInMemoryStore()
	categorySvc := category.NewService(idGenerator, categoryRepo, logger)
	NewCategoryRoutes(s, categorySvc)

	genreRepo := genre.NewInMemoryStore()
	genreSvc := genre.NewService(idGenerator, genreRepo, logger)
	NewGenreRoutes(s, genreSvc)

	castMemberRepo := inmem.NewCastMemberRepository()
	castMemberSvc := cast_member.NewService(idGenerator, castMemberRepo, logger)
	NewCastMemberRoutes(s, castMemberSvc)

	videoRepo := video.NewInMemoryStore()
	videoSvc := video.NewService(idGenerator, videoRepo, logger)
	NewVideoRoutes(s, videoSvc)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//s.logger.Info(r.Method, r.URL.Path)
	s.router.ServeHTTP(w, r)
}

func (s *server) AddRoute(method string, pattern string, handler http.Handler) {
	s.router.Handler(method, pattern, handler)
}

func (s *server) errBadRequest(w http.ResponseWriter, err error) {
	//s.logger.Warn(err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (s *server) errInternalServer(w http.ResponseWriter, err error) {
	//s.logger.Error(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *server) errNotFound(w http.ResponseWriter, err error) {
	//s.logger.Info(err)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *server) errUnprocessableEntity(w http.ResponseWriter, err error) {
	//s.logger.Warn(err)
	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

func (s *server) errStatusConflict(w http.ResponseWriter, err error) {
	//s.logger.Warn(err)
	http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
}
