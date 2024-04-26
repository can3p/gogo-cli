/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/can3p/go-scarf/scaffolder"
	gogoTemplate "github.com/can3p/gogo-cli/template"
	"github.com/google/uuid"

	"github.com/spf13/cobra"
)

func generateCommand() *cobra.Command {
	var test bool

	out := &cobra.Command{
		Use:   "generate <projectname>",
		Short: "generate a new gogo project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// All this is just to hack our way to a working project
			projectName := args[0]

			email, err := cmd.Flags().GetString("email")

			if err != nil {
				return err
			}
			if email == "" {
				return fmt.Errorf("email arg is required")
			}

			repo, err := cmd.Flags().GetString("repo")

			if err != nil {
				return err
			}
			if repo == "" {
				return fmt.Errorf("repo arg is required")
			}

			testEmail, err := cmd.Flags().GetString("testemail")

			if err != nil {
				return err
			}
			if testEmail == "" {
				return fmt.Errorf("testemail arg is required")
			}

			out, err := cmd.Flags().GetString("out")

			if err != nil {
				return err
			}

			if out == "" {
				path, err := os.Getwd()
				if err != nil {
					return err
				}

				out = path
			}

			out, err = filepath.Abs(out)

			if err != nil {
				return err
			}

			s := scaffolder.New().WithFileFilter(func(s string, d fs.DirEntry) scaffolder.FileFilterResult {
				if s == "template.go" {
					return scaffolder.FileFilterSkip
				}

				return scaffolder.FileFilterAccept
			}).WithFuncMap(map[string]any{
				"uuid": func() string {
					return uuid.NewString()
				}})

			if !test {
				s = s.WithProcessor(scaffolder.FSProcessor(out))
			}

			testEmailParts := strings.Split(testEmail, "@")

			return s.Scaffold(gogoTemplate.Template, scaffolder.ScaffoldData{
				"ProjectName":   projectName,
				"ProjectRepo":   repo,
				"ProjectEmail":  email,
				"TestEmailHead": testEmailParts[0],
				"TestEmailTail": testEmail[1],
			})
		},
	}

	out.Flags().String("email", "", "email to use as from address in the emails, mailjet should like it")
	out.Flags().String("repo", "", "project repo url like github.com/can3p/gogo-cli for the import names")
	out.Flags().String("testemail", "", "and address that you would use for testing. Only this address will be allowed to have a + sign for sign ups testing")
	out.Flags().String("out", "", "project will be generated in this folder if specified")
	out.Flags().BoolVarP(&test, "test", "", false, "Do not write anything, write everything to stdout")

	return out
}

func init() {
	rootCmd.AddCommand(generateCommand())
}
