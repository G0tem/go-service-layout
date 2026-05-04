package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/apple/dto"
	"github.com/G0tem/go-service-layout/internal/apple/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/google/uuid"
)

func (u *UseCase) CreateApple(ctx context.Context, input dto.CreateAppleInput) (dto.CreateAppleOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase CreateApple")
	defer span.End()

	var output dto.CreateAppleOutput

	a := entity.Apple{
		ID:     uuid.New(),
		Name:   input.Name,
		Status: entity.StatusNew,
	}

	err := u.postgres.CreateApple(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreateApple: %w", err)
	}

	err = u.redis.PutApple(ctx, a)
	if err != nil {
		return output, fmt.Errorf("u.redis.PutApple: %w", err)
	}

	event := entity.CreateEvent{
		ID:   a.ID,
		Name: input.Name,
	}

	err = u.kafka.CreateEvent(ctx, event)
	if err != nil {
		return output, fmt.Errorf("u.kafka.CreateEvent: %w", err)
	}

	output.ID = a.ID

	return output, nil
}
