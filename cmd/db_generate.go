package cmd

import (
	"fmt"
	"os"

	"github.com/can3p/gogo-cli/pkg/config"
	"github.com/can3p/gogo-cli/pkg/dbgen"
	"github.com/spf13/cobra"
)

func dbGenerateCommand() *cobra.Command {
	var projectDir string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate sqlboiler ORM code from database schema",
		Long: `Generate sqlboiler ORM code by:
1. Starting a PostgreSQL container
2. Applying all migrations from clean state
3. Running sqlboiler to generate Go code

Requires a gogo.toml configuration file in the project directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectDir == "" {
				var err error
				projectDir, err = os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
			}

			cfg, err := config.Load(projectDir)
			if err != nil {
				return err
			}

			generator := dbgen.NewGenerator(projectDir, cfg)
			return generator.Run()
		},
	}

	cmd.Flags().StringVarP(&projectDir, "dir", "d", "", "Project directory (default: current directory)")

	return cmd
}
