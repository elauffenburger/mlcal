package calcache

import (
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
)

type InMemoryCalendarCache struct {
	calendar *mlcal.Calendar
}

func NewInMemoryCache() *InMemoryCalendarCache {
	return &InMemoryCalendarCache{}
}

func (s *InMemoryCalendarCache) Get() (*mlcal.Calendar, error) {
	return s.calendar, nil
}

func (s *InMemoryCalendarCache) GetICS() (string, error) {
	if s.calendar == nil {
		return "", nil
	}

	return s.calendar.ToICS().Serialize(), nil
}

func (s *InMemoryCalendarCache) Set(cal *mlcal.Calendar) error {
	s.calendar = cal

	return nil
}
