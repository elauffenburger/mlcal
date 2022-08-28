package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar"
	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar/calcache"
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

const (
	flagEmail           string = "email"
	flagPassword        string = "password"
	flagCourseID        string = "course"
	flagRefreshInterval string = "refresh"
	flagRedisAddr       string = "redis-addr"
)

func main() {
	cmd := cobra.Command{
		Use: "api",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set up our server.
			srv := gin.Default()

			email, _ := cmd.Flags().GetString(flagEmail)
			password, _ := cmd.Flags().GetString(flagPassword)
			courseID, _ := cmd.Flags().GetString(flagCourseID)
			redisAddr, _ := cmd.Flags().GetString(flagRedisAddr)

			var refreshInterval time.Duration
			{
				str, _ := cmd.Flags().GetString(flagRefreshInterval)
				duration, err := time.ParseDuration(str)
				if err != nil {
					return err
				}

				refreshInterval = duration
			}

			// Add a health endpoint.
			srv.GET("/healthz", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			// Add the calendar endpoints.
			{
				logger := log.New(os.Stdout, "[calendar] ", log.LstdFlags)

				// Create the calendar cache.
				var calCache calendar.Cache
				if redisAddr != "" {
					logger.Print("using redis cache")

					cache, err := calendar.CreateRedisCache(&redis.Options{Addr: redisAddr})
					if err != nil {
						return err
					}

					calCache = cache
				} else {
					logger.Print("using in-memory cache")

					calCache = calcache.NewInMemoryCache()
				}

				client, err := mlcal.NewClient(email, password, courseID)
				if err != nil {
					return err
				}

				// Create a calendar refresher and set it up to refresh on an interval.
				calRefresher := calendar.NewAutoRefresher(logger, client, calCache)
				go calRefresher.RefreshOnInterval(refreshInterval, log.New(os.Stdout, "[calendar-refresher]: ", log.LstdFlags))

				if err := addCalendarEndpoints(srv, calRefresher); err != nil {
					return err
				}
			}

			// Start server.
			return srv.Run()
		},
	}

	cmd.Flags().String(flagEmail, "", "the email used to log in")
	cmd.MarkFlagRequired(flagEmail)

	cmd.Flags().String(flagPassword, "", "the password used to log in")
	cmd.MarkFlagRequired(flagPassword)

	cmd.Flags().String(flagCourseID, "", "the course ID to get cal for")
	cmd.MarkFlagRequired(flagCourseID)

	cmd.Flags().String(flagRefreshInterval, "24h", "the calendar refresh interval in time.ParseDuration format")

	cmd.Flags().String(flagRedisAddr, "", "the connection string to use for redis. If not provided, an in-memory cache wil be used instead")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error running mlcal api: %s\n", err)
		os.Exit(1)
	}
}
