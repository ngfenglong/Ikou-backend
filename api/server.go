package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/ngfenglong/ikou-backend/api/config"
	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/routes"
	"github.com/ngfenglong/ikou-backend/api/store"
)

type Application struct {
	Config   config.Config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Version  string
	Store    *store.Store

	PlaceController      *controllers.PlaceController
	AuthController       *controllers.AuthController
	CodestableController *controllers.CodestableController
}

type Server struct {
	Router chi.Router
	App    *Application
}

func NewServer(app *Application) *Server {
	server := &Server{
		Router: chi.NewRouter(),
		App:    app,
	}

	app.PlaceController = controllers.NewPlaceController(app.Store)
	app.CodestableController = controllers.NewCodestableController(app.Store)
	app.AuthController = controllers.NewAuthController(app.Store)

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TOKEN"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.Router.Mount("/api/places", routes.PlaceRoutes(s.App.Store))
	s.Router.Mount("/api/common", routes.CodestableRoutes(s.App.Store))
	s.Router.Mount("/api/auth", routes.AuthRoutes(s.App.Store))

}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.App.Config.Port),
		Handler:           s.Router,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	s.App.InfoLog.Printf("Backend is listening on port %d", s.App.Config.Port)

	return srv.ListenAndServe()
}
