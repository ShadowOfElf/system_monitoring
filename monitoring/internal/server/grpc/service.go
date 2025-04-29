package internal_grpc

import (
	"context"
	"sync"

	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	pb "github.com/ShadowOfElf/system_monitoring/pkg"
)

type ServiceGRPC struct {
	pb.UnimplementedMonitoringServer
	mu  sync.Mutex
	app *app.App
}

func NewGRPCService(app *app.App) *ServiceGRPC {
	return &ServiceGRPC{
		app: app,
	}
}

func (s *ServiceGRPC) GetStatisticProto(ctx context.Context, req *pb.GetStatistic) (*pb.StatisticResponse, error) {
	interval := req.StatsInterval
	// TODO добавить конкурентность
	statistic := s.app.GetStatistic(int(interval))
	return statToProtoStat(statistic), nil
}

func statToProtoStat(stat resources.Statistic) *pb.StatisticResponse {
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
