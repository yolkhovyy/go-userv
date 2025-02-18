package notifier

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
)

type Listener struct {
	*pq.Listener
}

func Connect(config postgres.Config) (*Listener, error) {
	// TODO: fix sslmode
	connString := "postgres://" +
		config.Username + ":" +
		config.Password + "@" +
		net.JoinHostPort(config.Host, strconv.Itoa(config.Port)) + "/" +
		config.Database + "?sslmode=disable"

	event := pq.ListenerEventDisconnected

	// The postgres listener starts asynchronously and
	// takes some time to connect to the databasse.
	// Using sync.WaitGroup to synchronize.
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	// TODO: make these configurable.
	const (
		minReconnectInterval = 5 * time.Second
		maxReconnectInterval = 30 * time.Second
	)

	var err error

	listener := pq.NewListener(connString,
		minReconnectInterval, maxReconnectInterval,
		func(levent pq.ListenerEventType, lerr error) {
			event = levent
			err = lerr

			waitGroup.Done()
		})

	// Wait until connected.
	waitGroup.Wait()

	if err != nil {
		return nil, fmt.Errorf("connect listener: %w", err)
	}

	if event != pq.ListenerEventConnected {
		return nil, fmt.Errorf("connect listener: %w %d", ErrListenerConnectFailed, event)
	}

	if listener == nil {
		return nil, fmt.Errorf("connect listener: %w", ErrNilListener)
	}

	err = listener.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &Listener{listener}, nil
}
