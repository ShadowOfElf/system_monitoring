package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ShadowOfElf/system_monitoring/configs"
	grpcclient "github.com/ShadowOfElf/system_monitoring/internal/client/grpc"
	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/olekukonko/tablewriter"
)

var (
	configString    string
	frequencyString string
	periodString    string
)

func init() {
	flag.StringVar(&configString, "config", "test.toml", "Path to configuration file")
	flag.StringVar(&frequencyString, "frequency", "10", "Frequency of requesting statistics")
	flag.StringVar(&periodString, "period", "10", "Period for which the statistics will be averaged")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config := configs.NewConfig(configString)
	logg := logger.New(config.Logger.Level)

	grpcService := grpcclient.NewGRPClient(logg, config.GRPC)

	tick, err := strconv.Atoi(frequencyString)
	if err != nil {
		logg.Error("Incorrect frequency")
	}

	period, err := strconv.Atoi(periodString)
	if err != nil {
		logg.Error("Incorrect period")
	}

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	client, err := grpcService.Start()
	if err != nil {
		logg.Error(fmt.Sprintf("failed create grpc client: %s", err))
		return
	}
	logg.Info("Start Client...")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Load", "CPU", "Disk", "Net"})
	table.SetAutoFormatHeaders(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	var lastLoad, lastCPU float32
	var lastDisk map[string]float32
	var lastNet map[string]int64

	updateTable := func(load, cpu float32, disk map[string]float32, net map[string]int64) {
		lastLoad, lastCPU, lastDisk, lastNet = load, cpu, disk, net

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("\033[2J\033[H")
			table.ClearRows()

			loadStr := ""
			if lastLoad >= 0 {
				loadStr = fmt.Sprintf("%.2f", lastLoad)
			}
			diskStr := mapToStr(lastDisk, func(s float32) string {
				return fmt.Sprintf("%0.2f%%", s)
			})
			netStr := mapToStr(lastNet, func(s int64) string {
				return fmt.Sprintf("%v", s)
			})

			table.Append([]string{
				loadStr,
				fmt.Sprintf("%.2f", lastCPU),
				diskStr,
				netStr,
			})
			table.Render()
		}
	}

	for {
		select {
		case <-ctx.Done():
			logg.Info("Stop Client...")
			return
		case <-ticker.C:
			stat, err := client.GetStatisticProto(ctx, grpcService.GetStatistic(period))
			if err != nil {
				logg.Error("Error with get stat:" + err.Error())
				continue
			}
			select {
			case <-ctx.Done():
				logg.Info("Stop Client...")
				return
			default:
				updateTable(stat.Statistic.Load, stat.Statistic.Cpu, stat.Statistic.Disk, stat.Statistic.Net)
			}
		}
	}
}

func mapToStr[V any](in map[string]V, f func(s V) string) string {
	keys := make([]string, 0, len(in))
	for name := range in {
		keys = append(keys, name)
	}
	sort.Strings(keys) // Сортируем ключи

	var sb strings.Builder
	for i, k := range keys {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(f(in[k]))
		if i < len(keys)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
