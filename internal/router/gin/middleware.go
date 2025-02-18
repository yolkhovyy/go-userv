package gin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger middleware logs HTTP requests.
func Logger() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		start := time.Now()

		log.Info().
			Str("method", gctx.Request.Method).
			Str("path", gctx.Request.URL.Path).
			Msg("http request begin")

		gctx.Next()

		log.Info().
			Str("method", gctx.Request.Method).
			Str("path", gctx.Request.URL.Path).
			Int("status", gctx.Writer.Status()).
			Dur("duration(ms)", time.Since(start).Round(time.Millisecond)).
			Msg("http request end")
	}
}
