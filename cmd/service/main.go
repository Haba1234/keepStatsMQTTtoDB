package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Haba1234/keepStatsMQTTtoDB/internal/app"
	"github.com/Haba1234/keepStatsMQTTtoDB/internal/clientmqtt"
	"github.com/Haba1234/keepStatsMQTTtoDB/internal/config"
	"github.com/Haba1234/keepStatsMQTTtoDB/internal/logger"
	"github.com/Haba1234/keepStatsMQTTtoDB/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/conf.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("config error: %v", err)
		os.Exit(1)
	}

	log := logger.NewLogger()
	log.Info(cfg)

	client := clientmqtt.NewClient(
		log,
		app.ConvertConfigClientMQTT(cfg.MQTT),
		"sprut",
		app.ConvertConfigServerMQTT(cfg.Servers["sprut"]),
	)

	db := storage.NewStorage(log, app.ConvertConfigStorage(cfg.Storage))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Канал для передачи точек в БД.
	pointsCh := make(chan app.Point, 200)

	if err = db.Start(ctx, pointsCh); err != nil {
		log.Error("failed to start Storage service: ", err.Error())
		cancel()
	}

	if err = client.Start(ctx, pointsCh); err != nil {
		log.Error("failed to start MQTT service: ", err.Error())
		cancel()
	}

	<-ctx.Done()

	if err := db.Stop(); err != nil {
		log.Error("failed to stop Storage service: ", err.Error())
	}
	if err := client.Stop(); err != nil {
		log.Error("failed to stop MQTT service: ", err.Error())
	}

	close(pointsCh)

	log.Info("shutdown complete")
}
