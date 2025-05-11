package internalgrpc

import (
	"fmt"
	"net"

	"github.com/ShadowOfElf/system_monitoring/configs"
	"github.com/ShadowOfElf/system_monitoring/internal/app"
	pb "github.com/ShadowOfElf/system_monitoring/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerGRPC struct {
	conf    configs.GRPCConf
	app     *app.App
	service *ServiceGRPC
	server  *grpc.Server
}

func NewServerGRPC(app *app.App, conf configs.GRPCConf) *ServerGRPC {
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerLogRequestInterceptor(app.Logger)))
	return &ServerGRPC{
		conf:    conf,
		app:     app,
		service: NewGRPCService(app),
		server:  server,
	}
}

func (s *ServerGRPC) Start() error {
	lsn, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		return err
	}
	pb.RegisterMonitoringServer(s.server, s.service)
	reflection.Register(s.server)
	s.app.Logger.Info(fmt.Sprintf("GRPC server is started on: %s", s.conf.Addr))
	go func() {
		if err := s.server.Serve(lsn); err != nil {
			s.app.Logger.Error(fmt.Sprintf("GRPC server start error: %s", err))
		}
	}()
	return nil
}

func (s *ServerGRPC) Stop() error {
	s.app.Logger.Info("GRPC server has been stopped")
	s.server.Stop()
	return nil
}
