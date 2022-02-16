package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type icsCalFetcher struct {
	client mlcal.Client
}

func (f *icsCalFetcher) Fetch() (string, error) {
	cal, err := f.client.Get()
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch calendar")
	}

	return cal.ToICS().Serialize(), nil
}

func newIcsCalFetcher(client mlcal.Client) *icsCalFetcher {
	return &icsCalFetcher{client}
}

func makeGetCalendarHandler(calFetcher *icsCalFetcher, refreshInterval *time.Duration, log func(string, ...interface{})) gin.HandlerFunc {
	lastFetched := time.Now()
	calMtx := sync.Mutex{}

	// Grab the calendar.
	cal, err := calFetcher.Fetch()
	if err != nil {
		panic(errors.Wrapf(err, "error fetching calendar"))
	}

	if refreshInterval != nil {
		// Set up a goroutine to refresh the calendar on a timer.
		go func() {
			for {
				log("waiting %s to refresh calendar", refreshInterval)

				// Wait for the refresh interval.
				<-time.After(*refreshInterval)

				log("refreshing calendar...")

				// Refresh the calendar.
				newCal, err := calFetcher.Fetch()
				if err != nil {
					log("error refreshing calendar: %s", err)
					continue
				}

				// Update the local cal.
				calMtx.Lock()
				lastFetched = time.Now()
				cal = newCal
				calMtx.Unlock()

				log("refreshed calendar")
			}
		}()
	}

	return func(c *gin.Context) {
		c.Header("Content-Type", "text/calendar")
		c.Header("Last-Modified", lastFetched.Format(http.TimeFormat))

		calMtx.Lock()
		c.String(200, "%s", cal)
		calMtx.Unlock()
	}
}
