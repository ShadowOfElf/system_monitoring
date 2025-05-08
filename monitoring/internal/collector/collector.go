package collector

import (
	"context"

	"github.com/ShadowOfElf/system_monitoring/internal/resources"
)

type InterfaceCollector interface {
	Start(ctx context.Context, tick int)
	Stop()
	Collect() resources.Snapshot
}
