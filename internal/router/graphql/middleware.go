package graphql

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

func withLogging(next graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		start := time.Now()

		log.Info().
			Str("request", params.Info.FieldName).
			Msgf("graphql request begin")

		nxt, err := next(params)

		log.Info().
			Str("request", params.Info.FieldName).
			Dur("duration(ms)", time.Since(start).Round(time.Millisecond)).
			Msg("graphql request end")

		return nxt, err
	}
}
