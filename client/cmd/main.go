package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

	var lastLoad, lastCPU, lastDisk, lastNet float32

	updateTable := func(load, cpu, disk, net float32) {
		lastLoad, lastCPU, lastDisk, lastNet = load, cpu, disk, net

		fmt.Print("\033[2K\r\033[1A\033[2K\r\033[1A\033[2K\r\033[1A\033[2K\r")

		select {
		case <-ctx.Done():
			return
		default:
			table.ClearRows()
			table.Append([]string{
				fmt.Sprintf("%.2f", lastLoad),
				fmt.Sprintf("%.2f", lastCPU),
				fmt.Sprintf("%.2f", lastDisk),
				fmt.Sprintf("%.2f", lastNet),
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
