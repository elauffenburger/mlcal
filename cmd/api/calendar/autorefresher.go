package calendar

import (
	"log"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar/calcache"
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
)

// AutoRefresher refreshes a calendar on a given interval that can later be retrieved.
type AutoRefresher struct {
	logger *log.Logger

	getter calcache.Getter
	cache  calcache.Cache
}

func NewAutoRefresher(logger *log.Logger, getter calcache.Getter, cache calcache.Cache) *AutoRefresher {
	return &AutoRefresher{logger, getter, cache}
}

func (g *AutoRefresher) Refresh() error {
	g.logger.Printf("refreshing calendar...")

	// Refresh the calendar.
	cal, err := g.getter.Get()
	if err != nil {
		g.logger.Printf("error fetching calendar: %s", err)
		return err
	}

	err = g.cache.Set(cal)
	if err != nil {
		g.logger.Printf("error caching calendar: %s", err)
		return err
	}

	g.logger.Printf("refreshed calendar")
	return nil
}

func (g *AutoRefresher) RefreshOnInterval(refreshInterval time.Duration) {
	for {
		_ = g.Refresh()

		// Wait for the refresh interval.
		g.logger.Printf("waiting %s to refresh calendar", refreshInterval)
		<-time.After(refreshInterval)
	}
}

func (g *AutoRefresher) Get() (*mlcal.Calendar, error) {
	return g.cache.Get()
}

func (g *AutoRefresher) GetICS() (string, error) {
	cal, err := g.Get()
	if err != nil {
		return "", err
	}

	if cal == nil {
		return "", nil
	}

	return cal.ToICS().Serialize(), nil
}
