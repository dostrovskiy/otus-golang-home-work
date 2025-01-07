package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/server/http"
	storage "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/sql"
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

	storage, err := newEventStorage(config)
	if err != nil {
		log.Error("failed to create event storage:", err.Error())
		return
	}
	defer storage.Close()

	calendar := app.New(log, storage)

	server := internalhttp.NewServer(log, calendar)

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

func newEventStorage(config *Config) (storage.EventStorage, error) {
	switch config.DataSource.StorageType {
	case "memory":
		return memorystorage.New(), nil
	case "sql":
		sqlstorage := sqlstorage.New(config.DataSource.Dsn)
		if err := sqlstorage.Open(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to connect to storage: %s", err.Error())
		}
		return sqlstorage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", config.DataSource.StorageType)
	}
}
