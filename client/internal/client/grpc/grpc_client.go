package grpcclient

import (
	"github.com/ShadowOfElf/system_monitoring/configs"
	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	pb "github.com/ShadowOfElf/system_monitoring/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceGRPCClient struct {
	log  logger.LogInterface
	conf configs.GRPCConf
}

func NewGRPClient(log logger.LogInterface, conf configs.GRPCConf) *ServiceGRPCClient {
	return &ServiceGRPCClient{
		log:  log,
		conf: conf,
	}
}

func (s *ServiceGRPCClient) Start() (pb.MonitoringClient, error) {
	conn, err := grpc.NewClient(s.conf.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewMonitoringClient(conn)
	s.log.Info("Start GRPC Client on:" + s.conf.Addr)
	return client, nil
}

func (s *ServiceGRPCClient) GetStatistic(period int) *pb.GetStatistic {
	return &pb.GetStatistic{StatsInterval: int64(period)}
}
