package notifier

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-userv/internal/contract/server"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
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
