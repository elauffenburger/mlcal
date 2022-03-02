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
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

const (
	flagEmail           string = "email"
	flagPassword        string = "password"
	flagCourse          string = "course"
	flagRefreshInterval string = "refresh"
	flagRedisAddr       string = "redis-addr"
)

func main() {
	cmd := cobra.Command{
		Use: "api",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set up our server.
			r := gin.Default()

			// Add the calendar endpoints.
			{
				logger := log.New(os.Stdout, "[calendar] ", log.LstdFlags)

				// Create the calendar cache.
				var calCache calendar.Cache
				redisAddr := cmd.Flag(flagRedisAddr).Value.String()
				if redisAddr != "" {
					logger.Print("using redis cache")

					cache, err := calendar.CreateRedisCache(&redis.Options{Addr: redisAddr})
					if err != nil {
						return err
					}

					calCache = cache
				} else {
					logger.Print("use in-memory cache")

					calCache = calcache.NewInMemoryCache()
				}

				// Create an AutoRefresher that can fetch the calendar on an interval.
				calRefresher := calendar.NewAutoRefresher(
					logger,
					mlcal.MustNewClient(
						cmd.Flag(flagEmail).Value.String(),
						cmd.Flag(flagPassword).Value.String(),
						cmd.Flag(flagCourse).Value.String(),
					),
					calCache,
				)

				// Kick off the refresh and set up a refresh on an interval.
				calRefresher.Refresh()
				go calRefresher.RefreshOnInterval(
					mustParseDuration(cmd.Flag(flagRefreshInterval).Value.String()),
					log.New(os.Stdout, "[calendar-refresher]: ", log.LstdFlags),
				)

				r.GET("/calendar", resource.Handler(calendar.MakeGetCalendarResource(calRefresher)))
				r.GET("/calendar/refresh", calendar.MakeRefreshCalendarEndpoint(calRefresher))
			}

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

	cmd.Flags().String(
		flagRedisAddr,
		"",
		"the connection string to use for redis. If not provided, an in-memory cache wil be used instead",
	)

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
