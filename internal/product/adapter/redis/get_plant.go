package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/product/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (r *Redis) GetPlant(ctx context.Context, id uuid.UUID) (entity.Plant, error) {
	ctx, span := tracer.Start(ctx, "redis GetApple")
	defer span.End()

	var plant entity.Plant

	data, err := r.client.Get(ctx, id.String()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return plant, entity.ErrNotFound
		}

		return plant, fmt.Errorf("r.client.Get: %w", err)
	}

	err = json.Unmarshal(data, &plant)
	if err != nil {
		return plant, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return plant, nil
}
