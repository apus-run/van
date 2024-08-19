package new

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// Cmd represents the new command.
var Cmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Example: "van new <your-project-name>",
	Short:   "Create a new project.",
	Long:    "create a new project with van layout.",
	Run:     run,
}

var (
	repoURL string
	timeout string
)

func init() {
	timeout = "60s"
	Cmd.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	Cmd.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
}

func run(_ *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	name := ""
	if len(args) == 0 {
		err := huh.NewInput().
			Title("What is your project name?").
			Description("project name.").
			Prompt("🚚").
			Value(&name).
			Validate(func(name string) error {
				if name == "" {
					return errors.New("The project name cannot be empty!")
				}
				return nil
			}).Run()

		if err != nil {
			return
		}
	} else {
		name = args[0]
	}
	p := NewProject(name)
	done := make(chan error, 1)
	go func() {
		done <- p.New(ctx, wd, repoURL)
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: project creation timed out\033[m\n")
		} else {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create project(%s)\033[m\n", ctx.Err().Error())
		}
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create project(%s)\033[m\n", err.Error())
		}
	}

}
