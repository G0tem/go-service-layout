package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/google/uuid"

	"github.com/G0tem/go-service-layout/pkg/transaction"
)

func (u *UseCase) CreatePineApple(ctx context.Context, _ dto.CreatePineAppleInput) (dto.CreatePineAppleOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase AddBanana")
	defer span.End()

	var output dto.CreatePineAppleOutput

	ctx, err := transaction.Begin(ctx)
	if err != nil {
		return output, fmt.Errorf("transaction.Begin: %w", err)
	}

	defer transaction.Rollback(ctx)

	err = u.postgres.CreateApple(ctx, entity.Apple{})
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreateApple: %w", err)
	}

	err = u.postgres.CreatePineApple(ctx, entity.PineApple{})
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreatePineApple: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return output, fmt.Errorf("transaction.Commit: %w", err)
	}

	output.ID = uuid.New()

	return output, nil
}
