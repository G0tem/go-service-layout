package kafka_producer

import (
	"context"
	"fmt"

	"github.com/G0tem/go-service-layout/internal/apple/entity"
	"github.com/G0tem/go-service-layout/pkg/otel/tracer"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func (p *Producer) CreateEvent(ctx context.Context, e entity.CreateEvent) error {
	ctx, span := tracer.Start(ctx, "kafka_producer UpdateEvent")
	defer span.End()

	m := kafka.Message{
		Key:   []byte(e.ID.String()),
		Value: []byte(e.Name),
	}

	err := p.writer.WriteMessages(ctx, m)
	if err != nil {
		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	log.Debug().Msg("CreateEvent Kafka")

	return nil
}
