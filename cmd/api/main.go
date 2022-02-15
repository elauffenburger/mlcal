package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	flagEmail           string = "email"
	flagPassword        string = "password"
	flagCourse          string = "course"
	flagRefreshInterval string = "refresh"
)

func main() {
	cmd := cobra.Command{
		Use: "api",
		RunE: func(cmd *cobra.Command, args []string) error {
			calClient, err := mlcal.NewClient(
				cmd.Flag(flagEmail).Value.String(),
				cmd.Flag(flagPassword).Value.String(),
				cmd.Flag(flagCourse).Value.String(),
			)
			if err != nil {
				return errors.Wrap(err, "error creating musical literacy client")
			}

			calFetcher := newIcsCalFetcher(calClient)
			calRefreshInterval, err := time.ParseDuration(cmd.Flag(flagRefreshInterval).Value.String())
			if err != nil {
				return errors.Wrap(err, "error parsing cal refresh duration")
			}

			// Start server.
			r := gin.Default()
			r.GET(
				"/calendar",
				makeGetCalendarHandler(
					calFetcher,
					&calRefreshInterval,
					func(s string, i ...interface{}) {
						fmt.Printf(s, i...)
						fmt.Println()
					},
				),
			)

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

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error running mlcal api: %s\n", err)
		os.Exit(1)
	}
}
