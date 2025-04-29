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

type InterfaceCollector interface {
	Start(ctx context.Context, tick int)
	Stop()
	Collect() resources.Snapshot
}

type Collector struct {
	app    *app.App
	cancel context.CancelFunc
}

func NewCollectorLinux(app *app.App) InterfaceCollector {
	return &Collector{
		app:    app,
		cancel: nil,
	}
}

func (c *Collector) Start(ctx context.Context, tick int) {
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

func (c *Collector) Stop() {
	if c.cancel != nil {
		c.app.Logger.Debug("Остановка сборщика")
		c.cancel()
	}
}

func (c *Collector) Collect() resources.Snapshot {
	load, err := CollectLoad()
	if err != nil {
		c.app.Logger.Error("Ошибка в получении загрузки:" + err.Error())
	}
	return resources.Snapshot{
		Load: load,
	}
}

func CollectLoad() (float32, error) {
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
