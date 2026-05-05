package http_server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string `envconfig:"HTTP_PORT" default:"8080"`
}

type Server struct {
	server *http.Server
	notify chan error
}

func New(handler http.Handler, port string) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         net.JoinHostPort("", port),
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}

	go s.start()

	log.Info().Msg("HTTP Server started on port: " + port)

	return s
}

func (s *Server) start() {
	s.notify <- s.server.ListenAndServe()
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully shuts down the HTTP server with a timeout
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server...")
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
		return fmt.Errorf("server shutdown: %w", err)
	}
	log.Info().Msg("HTTP server shut down successfully")
	return nil
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("server - Close - s.server.Shutdown")
	}

	log.Info().Msg("HTTP Server closed")
}
