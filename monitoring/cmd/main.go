package main

import (
	"flag"
	"fmt"

	"github.com/ShadowOfElf/system_monitoring/configs"
	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/logger"
)

var configString string

func init() {
	flag.StringVar(&configString, "config", "./test.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	config := configs.NewConfig(configString)
	logg := logger.New(config.Logger.Level)
	logg.Info("APP Started")
	application := app.New(logg)
	fmt.Println(application)
}
