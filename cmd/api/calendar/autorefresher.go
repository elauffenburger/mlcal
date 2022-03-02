package calendar

import (
	"log"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
)

// AutoRefresher refreshes a calendar on a given interval that can later be retrieved.
type AutoRefresher struct {
	logger *log.Logger

	getter Getter
	cache  Cache

	refreshed chan struct{}
}

func NewAutoRefresher(logger *log.Logger, getter Getter, cache Cache) *AutoRefresher {
	return &AutoRefresher{logger, getter, cache, make(chan struct{})}
}

func (r *AutoRefresher) Refresh() error {
	r.logger.Printf("refreshing calendar...")

	// Refresh the calendar.
	cal, err := r.getter.Get()
	if err != nil {
		r.logger.Printf("error fetching calendar: %s", err)
		return err
	}

	err = r.cache.Set(cal)
	if err != nil {
		r.logger.Printf("error caching calendar: %s", err)
		return err
	}

	r.logger.Printf("refreshed calendar")
	select {
	case r.refreshed <- struct{}{}:
	default:
	}
	return nil
}

func (r *AutoRefresher) RefreshOnInterval(refreshInterval time.Duration, logger *log.Logger) {
	r.Refresh()

	for {
		// Wait for the refresh interval.
		logger.Printf("waiting %s to refresh calendar", refreshInterval)

		select {
		// Wait until the refresh inverval passes.
		case <-time.After(refreshInterval):

		// ...but if the calendar is refreshed, skip this scheduled refresh.
		case <-r.refreshed:
			logger.Print("calendar manually refreshed; resetting refresh timer.")
			continue
		}

		_ = r.Refresh()
	}
}

func (r *AutoRefresher) Get() (*mlcal.Calendar, error) {
	return r.cache.Get()
}

func (r *AutoRefresher) GetICS() (string, error) {
	cal, err := r.Get()
	if err != nil {
		return "", err
	}

	if cal == nil {
		return "", nil
	}

	return cal.ToICS().Serialize(), nil
}
