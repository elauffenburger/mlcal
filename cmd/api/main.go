package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar"
	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar/calcache"
	"github.com/elauffenburger/musical-literacy-cal/cmd/api/resource"
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const (
	flagEmail           string = "email"
	flagPassword        string = "password"
	flagCourse          string = "course"
	flagRefreshInterval string = "refresh"
	flagRedisUrl        string = "redis-url"
)

func main() {
	cmd := cobra.Command{
		Use: "api",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create an AutoRefresher that can fetch the calendar on an interval.
			// Uses the calendar client to re-fetch calendars.
			// Uses an in-memory cache to store and retrieve the calendar.
			autoCalRefresher := calendar.NewAutoRefresher(
				log.New(os.Stdout, "[calendar] ", log.LstdFlags),
				mlcal.MustNewClient(
					cmd.Flag(flagEmail).Value.String(),
					cmd.Flag(flagPassword).Value.String(),
					cmd.Flag(flagCourse).Value.String(),
				),
				calcache.NewInMemoryCache(),
			)

			// Kick off the refresh and set up a refresh on an interval.
			autoCalRefresher.Refresh()
			go autoCalRefresher.RefreshOnInterval(mustParseDuration(cmd.Flag(flagRefreshInterval).Value.String()))

			// Set up our server.
			r := gin.Default()
			r.GET("/calendar", resource.Handler(calendar.MakeGetCalendarResource(autoCalRefresher)))

			// Start server.
			return r.Run()
		},
	}

	cmd.Flags().String(flagEmail, "", "the email used to log in")
	cmd.MarkFlagRequired(flagEmail)

	cmd.Flags().String(flagPassword, "", "the password used to log in")
	cmd.MarkFlagRequired(flagPassword)

	cmd.Flags().String(flagCourse, "", "the course ID to get cal for")
	cmd.MarkFlagRequired(flagCourse)

	cmd.Flags().String(flagRefreshInterval, "24h", "the calendar refresh interval in time.ParseDuration format")

	cmd.Flags().String(flagRedisUrl, "", "the url to use for redis; if omitted, results will be persisted in memory")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error running mlcal api: %s\n", err)
		os.Exit(1)
	}
}

func mustParseDuration(str string) time.Duration {
	d, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}

	return d
}
