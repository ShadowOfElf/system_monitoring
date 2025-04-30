//go:build linux

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

type LinuxCollector struct {
	app    *app.App
	cancel context.CancelFunc
	enable resources.CollectorEnable
}

func NewCollector(app *app.App, enable resources.CollectorEnable) InterfaceCollector {
	return &LinuxCollector{
		app:    app,
		cancel: nil,
		enable: enable,
	}
}

func (c *LinuxCollector) Start(ctx context.Context, tick int) {
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

func (c *LinuxCollector) Stop() {
	if c.cancel != nil {
		c.app.Logger.Debug("Остановка сборщика")
		c.cancel()
	}
}

func (c *LinuxCollector) Collect() resources.Snapshot {
	var err error
	var load float32 = -1
	if c.enable.Load {
		load, err = CollectLoadLin()
	}

	if err != nil {
		c.app.Logger.Error("Ошибка в получении загрузки:" + err.Error())
	}
	return resources.Snapshot{
		Load: load,
	}
}

func CollectLoadLin() (float32, error) {
	cmd := exec.Command("bash", "-c", "top -b -n 1 | grep \"%Cpu(s)\" | cut -d',' -f4 | awk '{print $1}'")
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
