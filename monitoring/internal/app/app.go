package app

import "github.com/ShadowOfElf/system_monitoring/internal/logger"

type App struct {
	Logger logger.LogInterface
}

func New(logg logger.LogInterface) *App {
	return &App{
		Logger: logg,
	}
}
