package main

import (
	"context"
	"flag"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ShadowOfElf/system_monitoring/configs"
	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/collector"
	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	internal_grpc "github.com/ShadowOfElf/system_monitoring/internal/server/grpc"
	"github.com/ShadowOfElf/system_monitoring/internal/storage"
)

var configString string

func init() {
	flag.StringVar(&configString, "config", "test.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	wg := sync.WaitGroup{}

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	config := configs.NewConfig(configString)
	logg := logger.New(config.Logger.Level)
	stor := storage.NewStorage(config.MaxSize, config.RepeatRateSec, logg, config.Enable)
	logg.Info("APP Started")
	application := app.New(logg, stor)
	grpcServer := internal_grpc.NewServerGRPC(application, config.GRPC)
	collectorService := collector.NewCollector(application, config.Enable)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err := grpcServer.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	collectorService.Start(ctx, config.RepeatRateSec)
	defer collectorService.Stop()

	logg.Info("Monitor server is running...")
	if err := grpcServer.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
	}

	wg.Wait()
}
