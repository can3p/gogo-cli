package cmd

import (
	"github.com/spf13/cobra"
)

func dbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database related commands",
	}

	cmd.AddCommand(dbGenerateCommand())

	return cmd
}

func init() {
	rootCmd.AddCommand(dbCommand())
}
