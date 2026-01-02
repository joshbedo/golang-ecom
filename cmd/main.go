package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8000",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	// Add structured logging - can filter by level=INFO,ERROR etc. or other metadata
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
