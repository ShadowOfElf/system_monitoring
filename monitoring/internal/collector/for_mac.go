//go:build darwin

package collector

import (
	"bytes"
	"context"
	"math"
	"os/exec"
	"strconv"
	"time"

	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
)

type MacCollector struct {
	app    *app.App
	cancel context.CancelFunc
	enable resources.CollectorEnable
}

func NewCollector(app *app.App, enable resources.CollectorEnable) InterfaceCollector {
	return &MacCollector{
		app:    app,
		cancel: nil,
		enable: enable,
	}
}

func (c *MacCollector) Start(ctx context.Context, tick int) {
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

func (c *MacCollector) Stop() {
	if c.cancel != nil {
		c.app.Logger.Debug("Остановка сборщика")
		c.cancel()
	}
}

func (c *MacCollector) Collect() resources.Snapshot {
	var err error
	var load float32 = -1
	if c.enable.Load {
		load, err = CollectLoadMac()
	}
	if err != nil {
		c.app.Logger.Error("Ошибка в получении загрузки:" + err.Error())
	}
	return resources.Snapshot{
		Load: load,
	}
}

func CollectLoadMac() (float32, error) {
	cmd := exec.Command("bash", "-c", "top -l 1 | grep \"CPU usage\" | awk '{print $7}' | cut -d'%' -f1")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	output = bytes.TrimSpace(output)
	result, err := strconv.ParseFloat(string(output), 32)
	if err != nil {
		return 0, err
	}

	load := 100 - result

	// Округляем до двух знаков после запятой
	roundedLoad := math.Round(load*100) / 100

	return float32(roundedLoad), nil
}
