package kafka_reader

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Addr  []string `envconfig:"KAFKA_READER_ADDR" required:"true"`
	Group string   `envconfig:"KAFKA_READER_GROUP" required:"true"`
	Topic string   `envconfig:"KAFKA_READER_TOPIC" required:"true"`
}

type Reader struct {
	*kafka.Reader
}

func New(c Config) (*Reader, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Addr,
		GroupID:  c.Group,
		Topic:    c.Topic,
		MaxBytes: 10e6, // 10MB
	})

	return &Reader{Reader: r}, nil
}

// Shutdown gracefully shuts down the Kafka reader with a timeout
func (r *Reader) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down Kafka reader...")

	done := make(chan error, 1)
	go func() {
		done <- r.Reader.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error().Err(err).Msg("Kafka reader shutdown error")
			return fmt.Errorf("kafka reader shutdown: %w", err)
		}
		log.Info().Msg("Kafka reader shut down successfully")
		return nil
	case <-ctx.Done():
		log.Warn().Msg("Kafka reader shutdown timeout")
		return fmt.Errorf("shutdown timeout: %w", ctx.Err())
	}
}

func (r *Reader) Close() {
	err := r.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("kafka_reader.Close")
	}

	log.Info().Msg("Kafka reader closed")
}
