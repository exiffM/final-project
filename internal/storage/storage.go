package storage

import (
	types "final-project/internal/statistics"
	"sync"
)

type Storage struct {
	mutex sync.Mutex
	stats []types.Statistic
}

func NewStorage() *Storage {
	return &Storage{stats: make([]types.Statistic, 0, 100)}
}

func (s *Storage) Append(st types.Statistic) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stats = append(s.stats, st)
}

func (s *Storage) Len() int {
	return len(s.stats)
}

func (s *Storage) PullOut(n int) []types.Statistic {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Len() < n {
		return s.getLastN(s.Len())
	}
	return s.getLastN(n)
}

func (s *Storage) getLastN(n int) []types.Statistic {
	lastIndex := len(s.stats)
	return s.stats[lastIndex-n:]
}

func (s *Storage) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stats = make([]types.Statistic, 0, 100)
}
