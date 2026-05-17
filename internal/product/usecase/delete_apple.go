package usecase

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/dto"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) DeleteApple(ctx context.Context, input dto.DeleteAppleInput) error {
	ctx, span := tracer.Start(ctx, "usecase DeleteApple")
	defer span.End()

	log := log.Ctx(ctx)

	log.Info().
		Str("apple_id", input.ID.String()).
		Msg("DeleteApple started")

	err := u.postgres.DeleteApple(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("u.postgres.DeleteApple: %w", err)
	}

	log.Info().
		Str("apple_id", input.ID.String()).
		Msg("DeleteApple completed")

	return nil
}
