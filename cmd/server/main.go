package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/G0tem/go-service-layout/config"
	_ "github.com/G0tem/go-service-layout/docs"
	"github.com/G0tem/go-service-layout/internal/app"
	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/otel"
	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"
)

// @title Order Service API
// @version 1.0
// @description Microservice for order management with JWT auth, Postgres, Redis, RabbitMQ
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}" to authenticate
func main() {
	// Create root context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		cancel()
	}()

	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	err = otel.Init(ctx, c.OTEL)
	if err != nil {
		log.Error().Err(err).Msg("otel.Init")
	}

	defer otel.Close()

	err = app.Run(ctx, c)
	if err != nil {
		log.Fatal().Err(err).Msg("app.Run")
	}

	log.Info().Msg("App stopped!")
}
