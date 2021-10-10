package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/Haba1234/keepStatsMQTTtoDB/internal/app"
	"github.com/Haba1234/keepStatsMQTTtoDB/internal/clientmqtt"
	"github.com/Haba1234/keepStatsMQTTtoDB/internal/config"
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
		log.Fatalf("Config error: %v", err)
	}
	log.Println(cfg)

	client := clientmqtt.NewClient(app.ConvertConfigClientMQTT(cfg.MQTT), "sprut", app.ConvertConfigServerMQTT(cfg.Servers["sprut"]))
	db := storage.NewStorage(app.ConvertConfigStorage(cfg.Storage))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Канал для передачи точек в БД.
	pointsCh := make(chan app.Point, 200)

	if err = db.Start(ctx, pointsCh); err != nil {
		log.Println("failed to start Storage service: " + err.Error())
		cancel()
	}

	if err = client.Start(ctx, pointsCh); err != nil {
		log.Println("failed to start MQTT service: " + err.Error())
		cancel()
	}

	<-ctx.Done()

	if err := db.Stop(); err != nil {
		log.Println("failed to stop Storage service: " + err.Error())
	}
	if err := client.Stop(); err != nil {
		log.Println("failed to stop MQTT service: " + err.Error())
	}

	close(pointsCh)

	log.Println("shutdown complete")
}
