package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const (
	flagEmail    string = "email"
	flagPassword string = "password"
	flagCourse   string = "course"
)

func main() {
	cmd := cobra.Command{
		Use: "api",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load the cal on startup.
			icsCal, err := loadICSCal(
				cmd.Flag(flagEmail).Value.String(),
				cmd.Flag(flagPassword).Value.String(),
				cmd.Flag(flagCourse).Value.String(),
			)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Start server.
			r := gin.Default()
			r.GET("/calendar", makeGetCalendarHandler(icsCal))

			return r.Run()
		},
	}

	cmd.Flags().String(flagEmail, "", "the email used to log in")
	cmd.MarkFlagRequired(flagEmail)

	cmd.Flags().String(flagPassword, "", "the password used to log in")
	cmd.MarkFlagRequired(flagPassword)

	cmd.Flags().String(flagCourse, "", "the course ID to get cal for")
	cmd.MarkFlagRequired(flagCourse)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error running mlcal api: %s\n", err)
		os.Exit(1)
	}
}
