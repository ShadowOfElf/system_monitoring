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
	mu       sync.Mutex
}

func NewStorage(maxSize int, log logger.LogInterface) InterfaceStorage {
	elements := make([]resources.Snapshot, maxSize)
	return &Storage{
		log:      log,
		elements: elements,
		len:      0,
		maxSize:  maxSize,
	}
}

func (s *Storage) Add(el resources.Snapshot) {
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
	stat := resources.Statistic{
		Load:       1234,
		CPU:        2345,
		Disk:       3456,
		Net:        4567,
		TopTalkers: []resources.TopTalker{resources.TopTalker{ID: 1, Name: "Test", LoadNet: 1111}},
	} // Тут будем получать стату из обработчика или вычислять на месте

	return stat
}
