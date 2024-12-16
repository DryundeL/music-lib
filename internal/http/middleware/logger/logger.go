package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger) gin.HandlerFunc {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	return func(c *gin.Context) {
		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.Request.RemoteAddr),
			slog.String("user_agent", c.Request.UserAgent()),
		)

		// Запоминаем время начала запроса
		startTime := time.Now()

		// Обрабатываем запрос
		c.Next()

		// Получаем статус ответа
		status := c.Writer.Status()
		bytesWritten := c.Writer.Size()
		if bytesWritten < 0 {
			bytesWritten = 0
		}

		entry.Info("request completed",
			slog.Int("status", status),
			slog.Int("bytes_written", bytesWritten),
			slog.String("duration", time.Since(startTime).String()),
		)
	}
}
