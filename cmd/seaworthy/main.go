package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gookit/color"
	"github.com/oklog/run"
	"github.com/spf13/cobra"

	"github.com/cakehappens/seaworthy/pkg/clioptions"
	cmdverify "github.com/cakehappens/seaworthy/pkg/cmd/verify"
	"github.com/cakehappens/seaworthy/pkg/kubernetes"
	"github.com/cakehappens/seaworthy/pkg/util/templates"

	// registers health checks
	_ "github.com/cakehappens/seaworthy/pkg/kubernetes/health/install"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

const binaryName = "seaworthy"

// NewSeaworthyCommand returns the root command for the CLI
func NewSeaworthyCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: binaryName,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
				return nil
			}

			fmt.Printf("seaworthy %s, commit %s, built at %s by %s", version, commit, date, builtBy)

			fmt.Printf("running: %s\n", args)

			fmt.Fprintln(os.Stderr, color.New(color.FgGreen).Sprint("Done!"))
			return nil
		},
	}

	ioStreams := clioptions.IOStreams{In: in, Out: out, ErrOut: err}
	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				cmdverify.New(ioStreams, kubernetes.GetResources),
			},
		},
	}
	groups.Add(cmd)
	templates.ActsAsRootCommand(cmd, nil, groups...)

	return cmd
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	runGroup := run.Group{}
	{
		cancelInterrupt := make(chan struct{})
		runGroup.Add(
			createSignalWatcher(ctx, cancelInterrupt, cancel),
			func(error) {
				close(cancelInterrupt)
			})
	}

	runGroup.Add(func() error {
		rootC := NewSeaworthyCommand(os.Stdin, os.Stdout, os.Stderr)
		rootC.SetArgs(os.Args[1:])
		return rootC.ExecuteContext(ctx)
	}, func(error) {
		cancel()
	})

	err := runGroup.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit reason: %s\n", err)
		os.Exit(1)
	}
}

// This function just sits and waits for ctrl-C
func createSignalWatcher(ctx context.Context, cancelInterruptChan <-chan struct{}, cancel context.CancelFunc) func() error {
	return func() error {
		c := make(chan os.Signal, 1)

		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			err := fmt.Errorf("received signal %s", sig)
			fmt.Fprintf(os.Stderr, "%s\n", err)
			signal.Stop(c)
			cancel()
			return err
		case <-ctx.Done():
			return nil
		case <-cancelInterruptChan:
			return nil
		}
	}
}
