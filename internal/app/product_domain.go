package app

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/product/adapter/kafka_producer"
	"github.com/G0tem/go-service-layout/internal/product/adapter/postgres"
	"github.com/G0tem/go-service-layout/internal/product/adapter/redis"
	"github.com/G0tem/go-service-layout/internal/product/controller/http_router"
	"github.com/G0tem/go-service-layout/internal/product/controller/kafka_consumer"
	"github.com/G0tem/go-service-layout/internal/product/usecase"
	"github.com/G0tem/go-service-layout/pkg/ratelimit"
)

func ProductDomain(d Dependencies, cfgRatelimit ratelimit.Config, ctx context.Context) {
	appleUseCase := usecase.New(
		postgres.New(d.Postgres.Pool),
		kafka_producer.New(d.KafkaWriter.Writer),
		redis.New(d.Redis.Client),
		d.TokenManager,
	)

	http_router.ProductRouter(d.RouterHTTP, appleUseCase, cfgRatelimit)

	go kafka_consumer.ProductConsumer(ctx, d.KafkaReader, appleUseCase)
}
