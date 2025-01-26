package main //nolint

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/logger"
	internalmessagebroker "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/messagebroker"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/notification"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s", err.Error())
		return
	}

	log := logger.New(config.Logger.Level)

	storage, err := app.NewEventStorage(config.DataSource.StorageType, config.DataSource.Dsn)
	if err != nil {
		log.Error("failed to create event storage: %s", err.Error())
		return
	}
	defer storage.Close()

	app := app.New(log, storage)

	producer := internalmessagebroker.NewProducer(log, config.Kafka.Brokers, config.Kafka.Topic)

	scheduler := notification.NewScheduler(log, app, producer)

	log.Info("Scheduler is starting\nConfig: %+v", config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if err := scheduler.Start(ctx, config.Scheduler.Interval); err != nil {
		log.Error("failed to start scheduler: %s", err.Error())
	}
}
//nolint