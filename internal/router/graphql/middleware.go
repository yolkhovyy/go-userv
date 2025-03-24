package graphql

import (
	"log/slog"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-otelw/pkg/tracew"
)

func withTelemetry(next graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (any, error) {
		var err error

		var nxt any

		start := time.Now()

		ctx, span := tracew.Start(params.Context, "graphql", params.Info.FieldName)
		defer func() { span.End(err) }()

		logger := slogw.DefaultLogger()

		logger.InfoContext(ctx, "graphql request begin",
			slog.String("field name", params.Info.FieldName),
		)

		params.Context = ctx
		nxt, err = next(params)

		logger.InfoContext(ctx, "graphql request end",
			slog.String("field name", params.Info.FieldName),
			slog.Duration("duration(ms)", time.Since(start).Round(time.Millisecond)),
		)

		return nxt, err
	}
}
