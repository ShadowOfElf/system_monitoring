package storage

import (
	"sync"

	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
)

// TODO: тут делаем интерфейс и реализацию через слайс с вытеснением принимающий размер слайса

type InterfaceStorage interface {
	Add(el resources.Snapshot)
	Len() int
	GetElements() []resources.Snapshot
	GetStatistic(interval int) resources.Statistic
}

type Storage struct {
	log      logger.LogInterface
	elements []resources.Snapshot
	len      int
	maxSize  int
	mu       sync.RWMutex
}

func NewStorage(maxSize int, log logger.LogInterface) InterfaceStorage {
	elements := make([]resources.Snapshot, 0, maxSize)
	return &Storage{
		log:      log,
		elements: elements,
		len:      0,
		maxSize:  maxSize,
	}
}

func (s *Storage) Add(el resources.Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.len == s.maxSize {
		copy(s.elements, s.elements[1:])
		s.elements = s.elements[:len(s.elements)-1]
		s.len -= 1
	}
	s.elements = append(s.elements, el)
	s.len += 1
}

func (s *Storage) Len() int {
	return s.len
}

func (s *Storage) GetElements() []resources.Snapshot {
	// TODO: Возможно удалить ближе к релизу, вроде ненужная фигня
	return s.elements
}

func (s *Storage) GetStatistic(interval int) resources.Statistic {
	//TODO не забыть или пересчитывать интервал в зависимости от repeat_rate или убрать repeat_rate
	s.mu.RLock()
	defer s.mu.RUnlock()
	lenElements := len(s.elements) - 1
	repeat := interval
	if repeat > lenElements {
		repeat = len(s.elements)
	}
	var load float32
	for i := lenElements; i >= lenElements-repeat; i-- {
		load += s.elements[i].Load
	}

	stat := resources.Statistic{
		Load:       load / float32(repeat),
		CPU:        2345,
		Disk:       3456,
		Net:        4567,
		TopTalkers: []resources.TopTalker{resources.TopTalker{ID: 1, Name: "Test", LoadNet: 1111}},
	}
	return stat
}
