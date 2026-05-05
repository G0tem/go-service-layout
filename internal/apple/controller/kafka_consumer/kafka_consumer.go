package kafka_consumer

import (
	"context"

	"github.com/G0tem/go-service-layout/internal/apple/usecase"
	"github.com/G0tem/go-service-layout/pkg/kafka_reader"
	"github.com/rs/zerolog/log"
)

func AppleConsumer(ctx context.Context, reader *kafka_reader.Reader, uc *usecase.UseCase) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Kafka consumer received shutdown signal")
			return
		default:
		}

		m, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("kafka_consumer.AppleConsumer: reader.FetchMessage")
			return
		}

		// UseCase call here

		log.Info().
			Str("topic", m.Topic).
			Int("partition", m.Partition).
			Int64("offset", m.Offset).
			Str("key", string(m.Key)).
			Str("value", string(m.Value)).
			Msg("kafka_consumer.AppleConsumer: reader.FetchMessage")

		if err = reader.CommitMessages(ctx, m); err != nil {
			log.Error().Err(err).Msg("kafka_consumer.AppleConsumer: reader.CommitMessages")
		}
	}
}
