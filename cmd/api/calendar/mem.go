package calendar

import (
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
)

type inMemoryCalendarCache struct {
	calendar *mlcal.Calendar
}

func NewInMemoryCache() Cache {
	return &inMemoryCalendarCache{}
}

func (s *inMemoryCalendarCache) Get() (*mlcal.Calendar, error) {
	return s.calendar, nil
}

func (s *inMemoryCalendarCache) GetICS() (string, error) {
	if s.calendar == nil {
		return "", nil
	}

	return s.calendar.ToICS().Serialize(), nil
}

func (s *inMemoryCalendarCache) Set(cal *mlcal.Calendar) error {
	s.calendar = cal

	return nil
}
