package storage

import (
	"testing"

	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	log := logger.New(logger.DebugLevel)
	maxSize := 10
	storage := Storage{
		log:        log,
		maxSize:    maxSize,
		repeatRate: 1,
		elements:   make([]resources.Snapshot, 0, maxSize),
		enable:     resources.CollectorEnable{Load: true},
		len:        0,
	}
	snapshot := resources.Snapshot{Load: 10}
	t.Run("storage add element", func(t *testing.T) {
		storage.Add(snapshot)
		storage.Add(snapshot)
		storage.Add(snapshot)
		elements := storage.GetElements()
		require.Equal(t, 3, len(elements))
	})

	t.Run("storage overflow", func(t *testing.T) {
		for i := 0; i < maxSize+1; i++ {
			storage.Add(snapshot)
		}
		elements := storage.GetElements()
		require.Equal(t, 10, len(elements))
	})

	t.Run("storage get statistic", func(t *testing.T) {
		storage.elements = make([]resources.Snapshot, 0, maxSize)
		storage.len = 0
		for i := 0; i < maxSize; i++ {
			snapshotTemp := resources.Snapshot{Load: float32(i * 10)}
			storage.Add(snapshotTemp)
		}
		statistic := storage.GetStatistic(10)
		require.Equal(t, float32(45), statistic.Load)

		statistic = storage.GetStatistic(5)
		require.Equal(t, float32(70), statistic.Load)
	})
}
