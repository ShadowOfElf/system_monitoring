//go:build linux

package collector

import (
	"bytes"
	"context"
	"math"
	"os/exec"
	"strconv"
	"strings"
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
		collectLoad: CollectLoad,
		collectCPU:  CollectCPU,
		collectDisk: CollectDisk,
		collectNet:  collectTCPStates,
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

func CollectLoad() (float32, error) {
	//
	cmd := exec.Command("bash", "-c", "cat /proc/loadavg | awk '{print $1}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	output = bytes.TrimSpace(output)

	result, err := strconv.ParseFloat(string(output), 32)
	if err != nil {
		return 0, err
	}
	return float32(result), nil
}

func CollectCPU() (float32, error) {
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

func CollectDisk() (map[string]float32, error) {
	cmd := exec.Command("bash", "-c", "df -hT")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	result := make(map[string]float32, len(lines)-1)
	for _, line := range lines[1:] { // Пропускаем заголовок
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		s := strings.TrimSuffix(fields[5], "%")
		value, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		result[fields[0]] = float32(value)
	}
	return result, nil
}

func collectTCPStates() (map[string]int64, error) {
	cmd := exec.Command("bash", "-c",
		`ss -tan | awk 'NR==1 {next} {state[$1]++} END {for(k in state) print k, state[k]}'`)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	states := make(map[string]int64)

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		state := strings.ReplaceAll(parts[0], "-", "_") // нормализуем названия (например, SYN-SENT → SYN_SENT)
		count, _ := strconv.Atoi(parts[1])
		states[state] = int64(count)
	}

	return states, nil
}
