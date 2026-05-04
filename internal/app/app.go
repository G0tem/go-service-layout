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

	waiting(httpServer)

	return nil
}

func waiting(httpServer *http_server.Server) {
	log.Info().Msg("App started!")

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)

	select {
	case i := <-wait:
		log.Info().Msg("App got signal: " + i.String())
	case err := <-httpServer.Notify():
		log.Error().Err(err).Msg("App got notify: httpServer.Notify")
	}

	log.Info().Msg("App is stopping...")
}
