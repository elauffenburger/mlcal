package calendar

import (
	"log"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
)

type AutoRefresherCache struct {
	logger *log.Logger

	g     Getter
	cache Cache
}

func NewAutoRefresher(logger *log.Logger, getter Getter, cache Cache) *AutoRefresherCache {
	return &AutoRefresherCache{logger, getter, cache}
}

func (g *AutoRefresherCache) Refresh() error {
	g.logger.Printf("refreshing calendar...")

	// Refresh the calendar.
	cal, err := g.g.Get()
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

func (g *AutoRefresherCache) RefreshOnInterval(refreshInterval time.Duration) {
	for {
		_ = g.Refresh()

		// Wait for the refresh interval.
		g.logger.Printf("waiting %s to refresh calendar", refreshInterval)
		<-time.After(refreshInterval)
	}
}

func (g *AutoRefresherCache) Get() (*mlcal.Calendar, error) {
	return g.cache.Get()
}

func (g *AutoRefresherCache) GetICS() (string, error) {
	cal, err := g.Get()
	if err != nil {
		return "", err
	}

	if cal == nil {
		return "", nil
	}

	return cal.ToICS().Serialize(), nil
}
