package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

// GinTraceLogger enriches request logs with trace context.
func GinTraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		spanContext := trace.SpanContextFromContext(c.Request.Context())
		traceID := spanContext.TraceID().String()
		spanID := spanContext.SpanID().String()

		reqLogger := log.With().
			Str("trace_id", traceID).
			Str("span_id", spanID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Logger()

		ctx := reqLogger.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		reqLogger.Info().
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Msg("http request handled")
	}
}

// Ctx returns logger from context or fallback global logger.
func Ctx(c *gin.Context) *zerolog.Logger {
	return log.Ctx(c.Request.Context())
}
