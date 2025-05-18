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

func (s *Storage) GetStatistic(interval int) resources.Statistic { //nolint
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

	var load, cpu float32
	disk := make(map[string]float32, lenElements)
	net := make(map[string]int64, lenElements)
	topT := make(map[string]int, 3)

	// если размеры будут очень большими вероятно стоит переделать на работу в параллельных горутинах
	for i := lenElements - 1; i >= lenElements-repeat; i-- {
		if s.enable.Load {
			load += s.elements[i].Load
		}

		if s.enable.CPU {
			cpu += s.elements[i].CPU
		}

		if s.enable.Disk {
			for name, value := range s.elements[i].Disk {
				disk[name] += value
			}
		}

		if s.enable.Net {
			for name, value := range s.elements[i].Net {
				net[name] += value
			}
		}

		if s.enable.TopTalkers {
			for _, talker := range s.elements[i].TopTalkers {
				topT[talker.Name] += talker.LoadNet
			}
		}
	}

	if s.enable.Load {
		load /= float32(repeat)
	} else {
		load = -1
	}
	if s.enable.CPU {
		cpu /= float32(repeat)
	} else {
		cpu = -1
	}

	if s.enable.Disk {
		for name, value := range disk {
			disk[name] = value / float32(repeat)
		}
	}

	if s.enable.Net {
		for name, value := range net {
			net[name] = value / int64(repeat)
		}
	}

	resultTalker := make([]resources.TopTalker, 0, 3)
	if s.enable.TopTalkers {
		for name, load := range topT {
			resultTalker = append(resultTalker, resources.TopTalker{Name: name, LoadNet: load / repeat})
		}
	}

	stat := resources.Statistic{
		Load:       load,
		CPU:        cpu,
		Disk:       disk,
		Net:        net,
		TopTalkers: resultTalker,
	}
	return stat
}
