package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/c95rt/context/config"

	"github.com/gorilla/mux"
	"github.com/joeshaw/envdecode"
	joonix "github.com/joonix/log"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func recoveryHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			debug.PrintStack()
			(&ResponseWriter{writer: w}).Error(http.StatusInternalServerError, "internal server error")
			return
		}
	}()
	next(w, r)
}

// AppHandlerFunc it's a custom definition for the http handlers.
type AppHandlerFunc func(*config.AppContext, *ResponseWriter, *http.Request)

// AppHandler it's a implementation of `http.Handler`, whose aim is
// wrap all hhtp handlers with app context.
type AppHandler struct {
	Context     *config.AppContext
	HandlerFunc AppHandlerFunc
}

func (a *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.HandlerFunc(a.Context, &ResponseWriter{writer: w}, r)
}

// Route Struct
type Route struct {
	Path    string
	Handler AppHandlerFunc
	// list of available HTTP methods
	Methods []string
	// it indicates that the url will be protected JWT
	IsProtected bool
}

// NewRouter instantiates a mux router
func NewRouter(ctx *config.AppContext, routes []*Route) *mux.Router {
	router := mux.NewRouter()
	for _, r := range routes {
		handler := &AppHandler{Context: ctx, HandlerFunc: r.Handler}
		if r.IsProtected {
			go router.Handle(r.Path, negroni.New(
				negroni.HandlerFunc(NewJWTMiddleware([]byte(ctx.Config.JWTSecret)).HandlerNext),
				negroni.Wrap(handler),
			)).Methods(r.Methods...)
		}
		if !r.IsProtected {
			go router.Handle(r.Path, handler).Methods(r.Methods...)
		}
	}
	return router
}

// UpServer ...
func UpServer(routes []*Route, plugins ...string) {
	log.SetFormatter(joonix.NewFormatter())
	var conf config.Configuration
	if err := envdecode.Decode(&conf); err != nil {
		fmt.Println(fmt.Errorf("could not load the app configuration: %v", err))
		log.Fatal(err)
	}

	context := &config.AppContext{
		Config: conf,
	}

	context.CloseGoRoutine = make(chan bool)

	server, err := createServer(context, routes)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Environment " + conf.Environment)
	log.Info("Listening on " + server.Addr)

	log.Fatal(server.ListenAndServe())
}

func createServer(context *config.AppContext, routes []*Route) (*http.Server, error) {
	n := negroni.New()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "HEAD"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "x-environment"},
	})
	n.Use(c)
	n.UseFunc(recoveryHandler)
	n.Use(negroni.HandlerFunc(LoggerRequest))
	n.Use(UserMiddleware())
	go n.UseHandler(NewRouter(context, routes))

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", context.Config.Port),
		ReadTimeout:  time.Duration(context.Config.Timeout) * time.Second,
		WriteTimeout: time.Duration(context.Config.Timeout) * time.Second,
		Handler:      n,
	}, nil
}
