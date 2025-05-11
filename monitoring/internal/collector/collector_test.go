package collector

import (
	"testing"

	"github.com/ShadowOfElf/system_monitoring/internal/app"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	"github.com/stretchr/testify/require"
)

type mockLogger struct{}

func (l *mockLogger) Debug(_ string) {}
func (l *mockLogger) Error(_ string) {}
func (l *mockLogger) Info(_ string)  {}
func (l *mockLogger) Warn(_ string)  {}

func TestCollector(t *testing.T) {
	t.Run("test_collect_from_mock", func(t *testing.T) {
		logger := &mockLogger{}
		application := &app.App{Logger: logger}

		collectorT := NewCollector(application, resources.CollectorEnable{Load: true})
		mockLoad := float32(75.5)
		collectorT.(*SCollector).collectLoad = func() (float32, error) {
			return mockLoad, nil
		}
		snapshot := collectorT.Collect()
		require.Equal(t, mockLoad, snapshot.Load)
	})
}
