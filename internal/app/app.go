package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/G0tem/go-service-layout/config"
	"github.com/G0tem/go-service-layout/pkg/http_server"
	jwt "github.com/G0tem/go-service-layout/pkg/jwt"
	"github.com/G0tem/go-service-layout/pkg/kafka_reader"
	"github.com/G0tem/go-service-layout/pkg/kafka_writer"
	"github.com/G0tem/go-service-layout/pkg/postgres"
	"github.com/G0tem/go-service-layout/pkg/redis"
	"github.com/G0tem/go-service-layout/pkg/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Dependencies struct {
	// Adapters
	Postgres    *postgres.Pool
	KafkaWriter *kafka_writer.Writer
	Redis       *redis.Client

	// Controllers
	RouterHTTP  *gin.Engine
	KafkaReader *kafka_reader.Reader

	// Auth
	TokenManager jwt.TokenManager
}

func Run(ctx context.Context, c config.Config) (err error) {
	var deps Dependencies

	// Adapters
	deps.Postgres, err = postgres.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer deps.Postgres.Close()

	deps.KafkaWriter, err = kafka_writer.New(c.KafkaWriter)
	if err != nil {
		return fmt.Errorf("kafka_writer.New: %w", err)
	}
	defer deps.KafkaWriter.Close()

	deps.Redis, err = redis.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}
	defer deps.Redis.Close()

	// Auth
	deps.TokenManager = jwt.NewManager(c.JWT.Secret, c.JWT.TTL)

	// Controllers
	deps.RouterHTTP = router.New()

	deps.KafkaReader, err = kafka_reader.New(c.KafkaReader)
	if err != nil {
		return fmt.Errorf("kafka_reader.New: %w", err)
	}
	defer deps.KafkaReader.Close()

	// Domains
	AppleDomain(deps)

	httpServer := http_server.New(deps.RouterHTTP, c.HTTP.Port)
	defer httpServer.Close()

    // Health checker
    healthChecker := healthcheck.New(c.HealthCheck)
    healthChecker.Start()
    defer healthChecker.Stop()

    // Rate limiting middleware
    if c.RateLimit.Enabled {
        deps.RouterHTTP.Use(ratelimit.RateLimiterMiddleware(c.RateLimit))
        log.Info().
            Int("rps", c.RateLimit.RequestsPerSecond).
            Int("burst", c.RateLimit.Burst).
            Msg("Rate limiting enabled")
        }

    waiting(httpServer, healthChecker, deps.KafkaReader)

	return nil
}

func waiting(httpServer *http_server.Server, healthChecker *healthcheck.HealthChecker, kafkaReader *kafka_reader.Reader) {
	log.Info().Msg("App started!")

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	shutdownCh := make(chan struct{})

	// Wait for shutdown signal
	wg.Add(1)
	go func() {
			defer wg.Done()
			wait := make(chan os.Signal, 1)
			signal.Notify(wait, os.Interrupt, syscall.SIGTERM)

			select {
			case i := <-wait:
					log.Info().Str("signal", i.String()).Msg("Received shutdown signal")
			case err := <-httpServer.Notify():
					if err != nil {
							log.Error().Err(err).Msg("HTTP server error")
					}
			}

			close(shutdownCh)
	}()

	// Wait for shutdown
	<-shutdownCh

	log.Info().Msg("App is stopping... Starting graceful shutdown")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown health checker first
	log.Info().Msg("Stopping health checker...")
	healthChecker.Stop()

	// Shutdown HTTP server
	log.Info().Msg("Shutting down HTTP server...")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	// Shutdown Kafka reader
	log.Info().Msg("Shutting down Kafka reader...")
	if err := kafkaReader.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Kafka reader shutdown error")
	}

	// Wait for all goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
			wg.Wait()
			close(done)
	}()

	select {
	case <-done:
			log.Info().Msg("All goroutines finished gracefully")
	case <-time.After(5 * time.Second):
			log.Warn().Msg("Timeout waiting for goroutines to finish")
	}

	log.Info().Msg("Graceful shutdown completed")
}
