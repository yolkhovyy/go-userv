package notifier

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-otelw/pkg/tracew"
	"github.com/yolkhovyy/go-userv/internal/contract/server"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
	"golang.org/x/sync/semaphore"
)

type Controller struct {
	pqListener  *pq.Listener
	kafkaWriter *kafka.Writer
}

//nolint:ireturn
func New(config postgres.Config, kafkaConfig Config) (server.Contract, error) {
	// Listener captures user changes in storage.
	listener, err := Connect(config)
	if err != nil {
		return nil, fmt.Errorf("storage listener: %w", err)
	}

	slogw.DefaultLogger().Debug("notifier connected to database")

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

	logger := slogw.DefaultLogger()

	err := c.pqListener.Listen(channel)
	if err != nil {
		return fmt.Errorf("notifier listen: %w", err)
	}

	defer func() {
		if err := c.kafkaWriter.Close(); err != nil {
			logger.ErrorContext(ctx, "notifier",
				slog.String("kafka writer close", err.Error()))
		}
	}()

	rateLimiter := semaphore.NewWeighted(rateLimit)

	for {
		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				logger.ErrorContext(ctx, "notifier",
					slog.String("context done", ctx.Err().Error()))

				return fmt.Errorf("notifier loop: %w", ctx.Err())
			}

			logger.DebugContext(ctx, "notifier",
				slog.String("listener exiting", ctx.Err().Error()))

			return nil

		case notification := <-c.pqListener.Notify:
			if notification == nil {
				continue
			}

			go func() {
				var err error

				ctx, span := tracew.Start(ctx, "listener", "notify")
				defer func() { span.End(err) }()

				logger := slogw.DefaultLogger()

				err = rateLimiter.Acquire(ctx, 1)
				if errors.Is(err, context.Canceled) {
					logger.WarnContext(ctx, "notifier listener exiting, context cancelled")

					return
				}
				defer rateLimiter.Release(1)

				span.AddEvent(fmt.Sprintf("writing topic %s %s", topic, notification.Extra))

				err = c.kafkaWriter.WriteMessages(ctx, kafka.Message{
					Topic: topic,
					Key:   []byte(key),
					Value: []byte(notification.Extra),
				})

				switch {
				case errors.Is(err, context.Canceled):
					logger.WarnContext(ctx, "notifier listener exiting, context cancelled")
				case err != nil:
					logger.ErrorContext(ctx, "notifier",
						slog.String("listener", err.Error()))
				default:
					logger.InfoContext(ctx, "notifier",
						slog.Any("notification extra", notification.Extra))
				}
			}()
		}
	}
}
