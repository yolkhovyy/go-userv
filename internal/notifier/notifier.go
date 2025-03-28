package notifier

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-otelw/pkg/tracew"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/semaphore"
)

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

				span.AddEvent("writing",
					trace.WithAttributes([]attribute.KeyValue{
						attribute.String("topic", topic),
						attribute.String("extra", notification.Extra),
					}...))

				headerCarrier := &KafkaHeaderCarrier{Headers: []kafka.Header{}}
				otel.GetTextMapPropagator().Inject(ctx, headerCarrier)

				err = c.kafkaWriter.WriteMessages(ctx, kafka.Message{
					Topic:   topic,
					Key:     []byte(key),
					Value:   []byte(notification.Extra),
					Headers: headerCarrier.Headers,
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
