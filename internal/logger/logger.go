// Package logger предоставляет инициализацию и использование логгера для всего приложения.
package logger

import (
	"net/http"

	"go.uber.org/zap"
)

// Log — глобальный логгер приложения.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует глобальный логгер с заданным уровнем логирования.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

// RequestLogger возвращает middleware, который логирует входящие HTTP-запросы.
func RequestLogger(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		h(w, r)
	})
}
