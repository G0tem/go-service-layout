package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/G0tem/go-service-layout/config"
	"github.com/G0tem/go-service-layout/pkg/healthcheck"
	"github.com/G0tem/go-service-layout/pkg/http_server"
	jwt "github.com/G0tem/go-service-layout/pkg/jwt"
	"github.com/G0tem/go-service-layout/pkg/kafka_reader"
	"github.com/G0tem/go-service-layout/pkg/kafka_writer"
	"github.com/G0tem/go-service-layout/pkg/postgres"
	"github.com/G0tem/go-service-layout/pkg/ratelimit"
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
	AppleDomain(deps, ctx)

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

	waiting(ctx, httpServer, healthChecker, deps.KafkaReader)

	return nil
}

func waiting(ctx context.Context, httpServer *http_server.Server, healthChecker *healthcheck.HealthChecker, kafkaReader *kafka_reader.Reader) {
	log.Info().Msg("App started!")

	// Wait for context cancellation (signal received)
	<-ctx.Done()

	log.Info().Msg("App is stopping... Starting graceful shutdown")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown components in order
	var wg sync.WaitGroup
	errCh := make(chan error, 3)

	// Shutdown health checker first
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Stopping health checker...")
		healthChecker.Stop()
	}()

	// Shutdown HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Shutting down HTTP server...")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			errCh <- fmt.Errorf("http server shutdown: %w", err)
			return
		}
	}()

	// Shutdown Kafka reader
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("Shutting down Kafka reader...")
		if err := kafkaReader.Shutdown(shutdownCtx); err != nil {
			errCh <- fmt.Errorf("kafka reader shutdown: %w", err)
			return
		}
	}()

	// Wait for all shutdown operations to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("All components shut down gracefully")
	case <-shutdownCtx.Done():
		log.Warn().Msg("Shutdown timeout exceeded")
	}

	// Log any errors
	close(errCh)
	for err := range errCh {
		log.Error().Err(err).Msg("Shutdown error")
	}

	log.Info().Msg("Graceful shutdown completed")
}
