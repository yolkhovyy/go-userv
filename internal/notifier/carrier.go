package notifier

import (
	"github.com/segmentio/kafka-go"
)

type KafkaHeaderCarrier struct {
	Headers []kafka.Header
}

func (k *KafkaHeaderCarrier) Get(key string) string {
	for _, header := range k.Headers {
		if header.Key == key {
			return string(header.Value)
		}
	}

	return ""
}

func (k *KafkaHeaderCarrier) Set(key string, value string) {
	for i, header := range k.Headers {
		if header.Key == key {
			k.Headers[i].Value = []byte(value)

			return
		}
	}

	k.Headers = append(k.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (k *KafkaHeaderCarrier) Keys() []string {
	keys := make([]string, len(k.Headers))
	for i, header := range k.Headers {
		keys[i] = header.Key
	}

	return keys
}
