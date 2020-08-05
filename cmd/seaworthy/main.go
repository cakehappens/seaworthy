package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cakehappens/seaworthy/pkg/clioptions"
	cmdverify "github.com/cakehappens/seaworthy/pkg/cmd/verify"
	"github.com/cakehappens/seaworthy/pkg/kubernetes"
	"github.com/cakehappens/seaworthy/pkg/util/templates"
	"github.com/gookit/color"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"syscall"

	// registers health checks
	_ "github.com/cakehappens/seaworthy/pkg/kubernetes/health/install"
)

const binaryName = "seaworthy"

func NewSeaworthyCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: binaryName,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				runHelp(cmd, args)
				return nil
			}

			fmt.Println("main")

			fmt.Printf("running: %s\n", args)


			fmt.Fprintln(os.Stderr, color.New(color.FgGreen).Sprint("Done!"))
			return nil
		},
	}

	client := kubernetes.NewKubectl()

	ioStreams := clioptions.IOStreams{In: in, Out: out, ErrOut: err}
	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				cmdverify.New(ioStreams, client),
			},
		},
	}
	groups.Add(cmd)
	templates.ActsAsRootCommand(cmd, nil, groups...)

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
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
	{
		runGroup.Add(func() error {
			rootC := NewSeaworthyCommand(os.Stdin, os.Stdout, os.Stderr)
			rootC.SetArgs(os.Args[1:])
			return rootC.ExecuteContext(ctx)
		}, func(error) {
			cancel()
		})
	}

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
			err := errors.New(fmt.Sprintf("received signal %s", sig))
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
