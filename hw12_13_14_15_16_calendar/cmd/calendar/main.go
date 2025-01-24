package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/logger"
	internalhttp "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v", err.Error())
		return
	}

	log := logger.New(config.Logger.Level)

	storage, err := app.NewEventStorage(config.DataSource.StorageType, config.DataSource.Dsn)
	if err != nil {
		log.Error("failed to create event storage:", err.Error())
		return
	}
	defer storage.Close()

	app := app.New(log, storage)

	server := internalhttp.NewServer(log, app)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop http server: %s", err.Error())
		}
	}()

	log.Info("Calendar is running.\nConfig: %+v", config)

	if err := server.Start(ctx, config.Server.Address); err != nil {
		log.Error("failed to start http server: %s", err.Error())
		cancel()
		return
	}
}
