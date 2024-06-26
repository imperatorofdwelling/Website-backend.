package http

import (
	"log"
	"time"

	"log/slog"
	"net/http"

	"github.com/imperatorofdwelling/Website-backend/internal/metrics"

	"github.com/imperatorofdwelling/Website-backend/pkg/repository/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerConfig struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	/*
		Idle timeout is a period of time during which
		the server or connection waits for any action from the client.
	*/
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}

func LoadConfig() (*ServerConfig, error) {
	//TODO load vars from .env
	return &ServerConfig{
		Addr:         "localhost:8000",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 5,
	}, nil
}

type Server struct {
	srv *http.Server
}

func New(cfg *ServerConfig, log *slog.Logger, repo postgres.LogRepository) *Server {
	srv := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      newRouter(log, repo),
	}
	return &Server{
		srv: srv,
	}
}

// Creating chi router
func newRouter(log *slog.Logger, repo postgres.LogRepository) http.Handler {
	r := chi.NewRouter()
	// There we need to write endpoints and middlewares

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)

	// We need db instance to work with it
	payment := metrics.NewPaymentHandler(log, repo)
	// TODO: names for endpoints
	r.Post(
		"/payment/create",
		payment.Payment)

	return r
}

func (s *Server) Run() {
	// Logger print need
	log.Println("Server start")
	if err := s.srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Disconnect() error {
	return s.srv.Close()
}
