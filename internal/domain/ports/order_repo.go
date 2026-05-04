package ports

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/domain/order"
)

type OrderRepository interface {
	Create(ctx context.Context, o order.Order) error
}
