package internalgrpc

import (
	"context"

	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	pb "github.com/ShadowOfElf/system_monitoring/pkg"
)

type ServiceGRPC struct {
	pb.UnimplementedMonitoringServer
	app *app.App
}

func NewGRPCService(app *app.App) *ServiceGRPC {
	return &ServiceGRPC{
		app: app,
	}
}

func (s *ServiceGRPC) GetStatisticProto(_ context.Context, req *pb.GetStatistic) (*pb.StatisticResponse, error) {
	interval := req.StatsInterval
	statistic := make(chan resources.Statistic, 1)

	go func() {
		statistic <- s.app.GetStatistic(int(interval))
	}()

	return statToProtoStat(statistic), nil
}

func statToProtoStat(statCh chan resources.Statistic) *pb.StatisticResponse {
	stat := <-statCh
	close(statCh)
	talkers := make([]*pb.TopTalker, 0, len(stat.TopTalkers))

	for _, talker := range stat.TopTalkers {
		protoTalker := pb.TopTalker{
			ID:   int64(talker.ID),
			Name: talker.Name,
			Load: talker.LoadNet,
		}
		talkers = append(talkers, &protoTalker)
	}
	protoStat := &pb.Statistic{
		Load:      stat.Load,
		Cpu:       stat.CPU,
		Disk:      stat.Disk,
		Net:       stat.Net,
		TopTalker: talkers,
	}
	return &pb.StatisticResponse{Statistic: protoStat}
}
