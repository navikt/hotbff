package hotbff

import (
	"log/slog"
	"os"
)

var (
	debugLog = os.Getenv("DEBUG_LOG") == "true"
)

func init() {
	level := slog.LevelInfo
	if debugLog {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
