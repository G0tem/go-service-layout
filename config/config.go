package config

import (
	"fmt"

	"github.com/G0tem/go-service-layout/pkg/http_server"
	"github.com/G0tem/go-service-layout/pkg/kafka_reader"
	"github.com/G0tem/go-service-layout/pkg/kafka_writer"
	"github.com/G0tem/go-service-layout/pkg/logger"
	"github.com/G0tem/go-service-layout/pkg/otel"
	"github.com/G0tem/go-service-layout/pkg/postgres"
	"github.com/G0tem/go-service-layout/pkg/redis"
	"github.com/G0tem/go-service-layout/pkg/sentry"
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App         App
	HTTP        http_server.Config
	Postgres    postgres.Config
	KafkaReader kafka_reader.Config
	KafkaWriter kafka_writer.Config
	Redis       redis.Config
	Logger      logger.Config
	Sentry      sentry.Config
	OTEL        otel.Config
}

func New() (Config, error) {
	var config Config

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	return config, fmt.Errorf("godotenv.Load: %w", err)
	// }

	err := envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
