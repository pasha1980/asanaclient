package main

import (
	"context"
	"github.com/pasha1980/asanaclient/internal/client"
	"github.com/pasha1980/asanaclient/internal/service"
	"github.com/pasha1980/asanaclient/internal/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM)

	asanaClient := client.NewAsanaClient()
	fileStorage := storage.NewStorage()

	slog.Info("Start extractor")
	go service.RunExtractor(ctx, service.ExtractorModeFiveMinutes, asanaClient, fileStorage)
	go service.RunExtractor(ctx, service.ExtractorModeThirtySeconds, asanaClient, fileStorage)
	slog.Info("Extractor is running")

	<-quit
	cancel()
}
