//go:build Windows

package collector

import (
	"context"
	"time"

	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
)

type WinCollector struct {
	app    *app.App
	cancel context.CancelFunc
	enable resources.CollectorEnable
}

func NewCollector(app *app.App, enable resources.CollectorEnable) InterfaceCollector {
	return &WinCollector{
		app:    app,
		cancel: nil,
		enable: enable,
	}
}

func (c *WinCollector) Start(ctx context.Context, tick int) {
	ctxCollector, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	go func() {
		for {
			select {
			case <-ctxCollector.Done():
				return
			case <-ticker.C:
				c.app.AddSnapshot(c.Collect())
			}
		}
	}()
	c.app.Logger.Debug("Запуск сборщика")
}

func (c *WinCollector) Stop() {
	if c.cancel != nil {
		c.app.Logger.Debug("Остановка сборщика")
		c.cancel()
	}
}

func (c *WinCollector) Collect() resources.Snapshot {
	var err error
	var load float32 = -1
	if c.enable.Load {
		load, err = CollectLoadWin()
	}
	if err != nil {
		c.app.Logger.Error("Ошибка в получении загрузки:" + err.Error())
	}
	return resources.Snapshot{
		Load: load,
	}
}

func CollectLoadWin() (float32, error) {
	// TODO реализовать
	return 0, nil
}
