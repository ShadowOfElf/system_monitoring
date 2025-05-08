package app

import (
	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	"github.com/ShadowOfElf/system_monitoring/internal/storage"
)

type App struct {
	Logger  logger.LogInterface
	Storage storage.InterfaceStorage
}

func New(logg logger.LogInterface, stor storage.InterfaceStorage) *App {
	return &App{
		Logger:  logg,
		Storage: stor,
	}
}

func (a *App) AddSnapshot(element resources.Snapshot) {
	a.Storage.Add(element)
}

func (a *App) GetStatistic(statsInterval int) resources.Statistic {
	return a.Storage.GetStatistic(statsInterval)
}
