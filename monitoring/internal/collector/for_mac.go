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

type SCollector struct {
	app         *app.App
	cancel      context.CancelFunc
	enable      resources.CollectorEnable
	collectLoad func() (float32, error)
	collectCPU  func() (float32, error)
	collectDisk func() (map[string]float32, error)
	collectNet  func() (map[string]int64, error)
}

func NewCollector(app *app.App, enable resources.CollectorEnable) InterfaceCollector {
	return &SCollector{
		app:         app,
		cancel:      nil,
		enable:      enable,
		collectLoad: CollectTEMPFloat,
		collectCPU:  CollectCPU,
		collectDisk: CollectTEMPDisk,
		collectNet:  CollectTEMPNet,
	}
}

func (c *SCollector) Start(ctx context.Context, tick int) {
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

func (c *SCollector) Stop() {
	if c.cancel != nil {
		c.app.Logger.Debug("Остановка сборщика")
		c.cancel()
	}
}

func (c *SCollector) Collect() resources.Snapshot {
	var err error
	var load float32 = -1
	var cpu float32 = -1
	var disk map[string]float32
	var net map[string]int64

	if c.enable.Load {
		load, err = c.collectLoad()
	}
	if err != nil {
		c.app.Logger.Error("Ошибка в получении загрузки:" + err.Error())
	}

	if c.enable.CPU {
		cpu, err = c.collectCPU()
	}
	if err != nil {
		c.app.Logger.Error("Ошибка в получении CPU:" + err.Error())
	}

	if c.enable.Disk {
		disk, err = c.collectDisk()
	}
	if err != nil {
		c.app.Logger.Error("Ошибка в получении Disk:" + err.Error())
	}

	if c.enable.Net {
		net, err = c.collectNet()
	}
	if err != nil {
		c.app.Logger.Error("Ошибка в получении Net:" + err.Error())
	}

	return resources.Snapshot{
		Load: load,
		CPU:  cpu,
		Disk: disk,
		Net:  net,
	}
}

func CollectTEMPFloat() (float32, error) {
	return -1, nil
}

func CollectTEMPNet() (map[string]int64, error) {
	return nil, nil
}

func CollectTEMPDisk() (map[string]float32, error) {
	return nil, nil
}

func CollectCPU() (float32, error) {
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
