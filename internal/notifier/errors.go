package notifier

import "errors"

var (
	ErrListenerConnectFailed = errors.New("listener connect failed")
	ErrNilListener           = errors.New("nil listener")
)
