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

func (r *Redis) GetApple(ctx context.Context, id uuid.UUID) (entity.Apple, error) {
	ctx, span := tracer.Start(ctx, "redis GetApple")
	defer span.End()

	var apple entity.Apple

	data, err := r.client.Get(ctx, id.String()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return apple, entity.ErrNotFound
		}

		return apple, fmt.Errorf("r.client.Get: %w", err)
	}

	err = json.Unmarshal(data, &apple)
	if err != nil {
		return apple, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return apple, nil
}
