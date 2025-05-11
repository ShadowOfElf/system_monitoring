package internalgrpc

import (
	"context"
	"testing"

	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	"github.com/ShadowOfElf/system_monitoring/internal/storage"
	pb "github.com/ShadowOfElf/system_monitoring/pkg"
	"github.com/stretchr/testify/require"
)

func TestServicesGRPC(t *testing.T) {
	logg := logger.New(logger.DebugLevel)
	enable := resources.CollectorEnable{
		Load: true,
	}
	store := storage.NewStorage(100, 1, logg, enable)
	application := app.New(logg, store)
	service := NewGRPCService(application)

	t.Run("test_grpc", func(t *testing.T) {
		req := &pb.GetStatistic{StatsInterval: 10}
		resp, err := service.GetStatisticProto(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
