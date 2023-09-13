package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func LogMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		baseLogMessage := logger.With(
			zap.String("path", string(c.Request().URI().Path())),
			zap.String("method", string(c.Request().Header.Method())),
		)

		requestQuery := c.Request().URI().QueryArgs().String()
		if requestQuery != "" {
			baseLogMessage.With(zap.String("query", requestQuery))
		}

		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start)

		if err != nil {
			baseLogMessage.Log(zap.ErrorLevel, err.Error(), zap.Duration("time", elapsed))
			return err
		}

		baseLogMessage.Log(zap.InfoLevel, "completed request", zap.Duration("time", elapsed))
		return nil
	}
}
