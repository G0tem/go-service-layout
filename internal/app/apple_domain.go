package app

import (
	"github.com/G0tem/go-service-layout/internal/apple/adapter/kafka_producer"
	"github.com/G0tem/go-service-layout/internal/apple/adapter/postgres"
	"github.com/G0tem/go-service-layout/internal/apple/adapter/redis"
	"github.com/G0tem/go-service-layout/internal/apple/controller/http_router"
	"github.com/G0tem/go-service-layout/internal/apple/controller/kafka_consumer"
	"github.com/G0tem/go-service-layout/internal/apple/usecase"
)

func AppleDomain(d Dependencies) {
	appleUseCase := usecase.New(
		postgres.New(d.Postgres.Pool),
		kafka_producer.New(d.KafkaWriter.Writer),
		redis.New(d.Redis.Client),
		d.TokenManager,
	)

	http_router.AppleRouter(d.RouterHTTP, appleUseCase)

	go kafka_consumer.AppleConsumer(d.KafkaReader, appleUseCase)
}
