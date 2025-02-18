package notifier

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/yolkhovyy/user/internal/contract/server"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
	"golang.org/x/sync/semaphore"
)

type Controller struct {
	pqListener  *pq.Listener
	kafkaWriter *kafka.Writer
}

//nolint:ireturn
func New(storageConfig storage.Config, kafkaConfig Config) (server.Contract, error) {
	// Listener captures user changes in storage.
	listener, err := Connect(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("storage listener: %w", err)
	}

	log.Debug().Msg("notifier connected to database")

	// TODO: make kafka params configurable.
	const (
		batchSize    = 1e6
		batchBytes   = 1e6
		batchTimeout = 100 * time.Millisecond
	)

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      kafkaConfig.Brokers,
		BatchTimeout: batchTimeout,
		BatchSize:    batchSize,
		BatchBytes:   batchBytes,
	})

	controller := &Controller{
		kafkaWriter: kafkaWriter,
		pqListener:  listener.Listener,
	}

	return controller, nil
}

//nolint:funlen,cyclop
func (c *Controller) Run(ctx context.Context) error {
	// TODO: Inspect hard-coded values.
	const (
		topic     = "postgres.public.users"
		key       = "user-event"
		channel   = "user_changes"
		rateLimit = 500
	)

	err := c.pqListener.Listen(channel)
	if err != nil {
		return fmt.Errorf("notifier listen: %w", err)
	}

	defer func() {
		if err := c.kafkaWriter.Close(); err != nil {
			log.Error().Err(err).Msg("kafka writer close")
		}
	}()

	rateLimiter := semaphore.NewWeighted(rateLimit)

	for {
		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				log.Error().Err(ctx.Err()).Msg("notifier")

				return fmt.Errorf("notifier loop: %w", ctx.Err())
			}

			log.Trace().Msg("notifier listener exiting")

			return nil

		case notification := <-c.pqListener.Notify:
			if notification == nil {
				continue
			}

			go func() {
				err := rateLimiter.Acquire(ctx, 1)
				if errors.Is(err, context.Canceled) {
					log.Warn().Msg("notifier listener exiting, context cancelled")

					return
				}
				defer rateLimiter.Release(1)

				err = c.kafkaWriter.WriteMessages(ctx, kafka.Message{
					Topic: topic,
					Key:   []byte(key),
					Value: []byte(notification.Extra),
				})

				switch {
				case errors.Is(err, context.Canceled):
					log.Warn().Msg("notifier listener exiting, context cancelled")
				case err != nil:
					log.Error().Err(err).Msg("notifier listener")
				default:
					log.Info().Msgf("notifier: %v", notification.Extra)
				}
			}()
		}
	}
}
