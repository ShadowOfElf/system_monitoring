package storage

import (
	"math"
	"sync"

	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
)

type InterfaceStorage interface {
	Add(el resources.Snapshot)
	Len() int
	GetElements() []resources.Snapshot
	GetStatistic(interval int) resources.Statistic
}

type Storage struct {
	log        logger.LogInterface
	elements   []resources.Snapshot
	len        int
	maxSize    int
	repeatRate int
	mu         sync.RWMutex
	enable     resources.CollectorEnable
}

func NewStorage(
	maxSize int, repeatRate int, log logger.LogInterface, enable resources.CollectorEnable,
) InterfaceStorage {
	elements := make([]resources.Snapshot, 0, maxSize)
	return &Storage{
		log:        log,
		elements:   elements,
		len:        0,
		maxSize:    maxSize,
		repeatRate: repeatRate,
		enable:     enable,
	}
}

func (s *Storage) Add(el resources.Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.len == s.maxSize {
		copy(s.elements, s.elements[1:])
		s.elements = s.elements[:len(s.elements)-1]
		s.len--
	}
	s.elements = append(s.elements, el)
	s.len++
}

func (s *Storage) Len() int {
	return s.len
}

func (s *Storage) GetElements() []resources.Snapshot {
	return s.elements
}

func (s *Storage) GetStatistic(interval int) resources.Statistic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	lenElements := len(s.elements)

	if lenElements < 0 {
		return resources.Statistic{}
	}

	repeat := int(math.Round(float64(interval) / float64(s.repeatRate)))
	if repeat > lenElements {
		repeat = len(s.elements)
	}

	var load float32
	if s.enable.Load {
		for i := lenElements - 1; i >= lenElements-repeat; i-- {
			load += s.elements[i].Load
		}
		load /= float32(repeat)
	} else {
		load = -1
	}

	stat := resources.Statistic{
		Load:       load,
		CPU:        2345,
		Disk:       3456,
		Net:        4567,
		TopTalkers: []resources.TopTalker{{ID: 1, Name: "Test", LoadNet: 1111}},
	}
	return stat
}
