package usecase

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/apple/entity"
	"github.com/google/uuid"
)

type Postgres interface {
	CreateApple(ctx context.Context, a entity.Apple) (err error)
	GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error)

	CreatePineApple(ctx context.Context, p entity.PineApple) (err error)
}

type Kafka interface {
	CreateEvent(ctx context.Context, e entity.CreateEvent) error
}

type Redis interface {
	GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error)
	PutApple(ctx context.Context, a entity.Apple) error
}

type UseCase struct {
	postgres Postgres
	kafka    Kafka
	redis    Redis
}

func New(p Postgres, k Kafka, r Redis) *UseCase {
	return &UseCase{
		postgres: p,
		kafka:    k,
		redis:    r,
	}
}
