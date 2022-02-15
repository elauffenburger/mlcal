package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/spf13/cobra"
)

const (
	flagEmail    string = "email"
	flagPassword string = "password"
	flagCourse   string = "course"
)

func main() {
	cmd := cobra.Command{
		Use: "mlcal",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Make a new Musical Literacy client.
			client, err := mlcal.NewClient(
				cmd.Flag(flagEmail).Value.String(),
				cmd.Flag(flagPassword).Value.String(),
				cmd.Flag(flagCourse).Value.String(),
			)
			if err != nil {
				return err
			}

			// Grab the calendar.
			cal, err := client.Get()
			if err != nil {
				return err
			}

			// Convert the calendar to an ICS format.
			icsCal := cal.ToICS()

			// Write the ICS to stdout.
			out := bufio.NewWriter(os.Stdout)
			icsCal.SerializeTo(out)

			err = out.Flush()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagEmail, "", "the email used to log in")
	cmd.MarkFlagRequired(flagEmail)

	cmd.Flags().String(flagPassword, "", "the password used to log in")
	cmd.MarkFlagRequired(flagPassword)

	cmd.Flags().String(flagCourse, "", "the course ID to get cal for")
	cmd.MarkFlagRequired(flagCourse)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error running mlcal: %s\n", err)
		os.Exit(1)
	}
}
