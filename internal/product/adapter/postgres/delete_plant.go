package postgres

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (p *Postgres) DeletePlant(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "postgres DeletePlant")
	defer span.End()

	dataset := goqu.Delete("plant").
		Where(goqu.C("id").Eq(id))

	sql, _, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("dataset.ToSQL: %w", err)
	}

	_, err = p.pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("p.pool.Exec: %w", err)
	}

	return nil
}
