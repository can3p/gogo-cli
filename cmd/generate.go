/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func generateCommand() *cobra.Command {
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

			testemail, err := cmd.Flags().GetString("testemail")

			if err != nil {
				return err
			}
			if testemail == "" {
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

			return generateProject(projectName, out, email, repo, testemail)
		},
	}

	out.Flags().String("email", "", "email to use as from address in the emails, mailjet should like it")
	out.Flags().String("repo", "", "project repo url like github.com/can3p/gogo-cli for the import names")
	out.Flags().String("testemail", "", "and address that you would use for testing. Only this address will be allowed to have a + sign for sign ups testing")
	out.Flags().String("out", "", "project will be generated in this folder if specified")

	return out
}

func generateProject(projectName, out, email, repo, testemail string) error {
	targetFolder := out + "/" + projectName

	if err := exec.Command("cp", "-R", "template", targetFolder).Run(); err != nil {
		return err
	}

	testemailParts := strings.Split(testemail, "@")
	from := []string{"projectname", "projectemail", "projectrepo", "testemailhead", "testemailtail"}
	to := []string{projectName, email, repo, testemailParts[0], testemailParts[1]}

	for idx := range from {
		field := "<" + from[idx] + ">"
		val := strings.ReplaceAll(to[idx], "/", "\\/")

		// find . -name '*.go' -exec sed -i 's/<reponame>/<projectrepo>/g' {} \;
		//args := []string{"find", out, "-type", "f", "-exec", fmt.Sprintf("sed -i 's/%s/%s/g' {}", field, val), ";"}
		args := []string{"find", out, "-type", "f"}
		log.Println("sed args:", strings.Join(args, " "))
		cmd := exec.Command(args[0], args[1:]...)

		if out, err := cmd.Output(); err != nil {
			return err
		} else {
			fnames := strings.Split(string(out), "\n")

			for _, fname := range fnames {
				if fname == "" {
					continue
				}

				args := []string{"sed", "-i", fmt.Sprintf("s/%s/%s/g", field, val), fname}
				log.Println("sed args:", strings.Join(args, " "))
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(generateCommand())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
