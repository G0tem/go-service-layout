package usecase

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	jwt "github.com/G0tem/go-service-layout/pkg/jwt"
	"github.com/google/uuid"
)

type Postgres interface {
	CreateApple(ctx context.Context, a entity.Apple) (err error)
	GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error)
	UpdateApple(ctx context.Context, a entity.Apple) (err error)
	DeleteApple(ctx context.Context, id uuid.UUID) error

	CreatePlant(ctx context.Context, p entity.Plant) (err error)
	GetPlant(ctx context.Context, id uuid.UUID) (entity.Plant, error)
	UpdatePlant(ctx context.Context, p entity.Plant) (err error)
	DeletePlant(ctx context.Context, id uuid.UUID) error
}

type Kafka interface {
	CreateEvent(ctx context.Context, e entity.CreateEvent) error
}

type Redis interface {
	GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error)
	PutApple(ctx context.Context, a entity.Apple) error
	// DeleteApple(ctx context.Context, id uuid.UUID) error

	GetPlant(ctx context.Context, id uuid.UUID) (entity.Plant, error)
	PutPlant(ctx context.Context, a entity.Plant) error
	// DeletePlant(ctx context.Context, id uuid.UUID) error
}

type UseCase struct {
	postgres     Postgres
	kafka        Kafka
	redis        Redis
	tokenManager jwt.TokenManager
}

func New(p Postgres, k Kafka, r Redis, tm jwt.TokenManager) *UseCase {
	return &UseCase{
		postgres:     p,
		kafka:        k,
		redis:        r,
		tokenManager: tm,
	}
}

func (u *UseCase) GetTokenManager() jwt.TokenManager {
	return u.tokenManager
}
