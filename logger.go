package hotbff

import (
	"log/slog"
	"os"
)

func init() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
