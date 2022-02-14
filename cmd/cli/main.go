package cli

import (
	"fmt"
	"os"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use: "mlcal",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := mlcal.NewClient(
				cmd.Flag("email").Value.String(),
				cmd.Flag("password").Value.String(),
				cmd.Flag("courseID").Value.String(),
			)

			cal, err := client.Get()
			if err != nil {
				return err
			}
		},
	}

	cmd.Flags().String("email", "", "the email to log in as")
	cmd.Flags().String("password", "", "the password to use to log in")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error running mlcal: %s", err)
		os.Exit(1)
	}
}
