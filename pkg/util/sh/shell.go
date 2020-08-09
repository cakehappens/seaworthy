package sh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// RunOptions is part of the functional API for Run
type RunOptions struct {
	Args   []string
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// RunOption is part of the functional API for Run
type RunOption func(o *RunOptions)

// Run provides a functional API for running a process
// Default Behaviors:
// stdout, stderr, stdin are sent to /dev/null
// all environment variables are passed to subprocess
// additionally, you may add vars to RunOptions.Env
//
// Returns the exitCode and error
// The error turned will only be non-nil if the process failed to start for some reason
func Run(ctx context.Context, cmd string, options ...RunOption) (int, error) {
	runOpts := &RunOptions{}

	runOpts.Env = make(map[string]string)

	for _, optFn := range options {
		optFn(runOpts)
	}

	cmdEnv := make([]string, len(runOpts.Env)+len(os.Environ()))
	cmdEnv = append(cmdEnv, os.Environ()...)

	for k, v := range runOpts.Env {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v))
	}

	c := exec.CommandContext(ctx, cmd, runOpts.Args...)
	c.Env = cmdEnv

	// fmt.Printf("cmd: %s %s\n", cmd, strings.Join(runOpts.Args, " "))

	c.Stdout = runOpts.Stdout
	c.Stderr = runOpts.Stderr
	c.Stdin = runOpts.Stdin

	err := c.Run()
	if _, ok := err.(*exec.ExitError); ok {
		return c.ProcessState.ExitCode(), nil
	}

	return c.ProcessState.ExitCode(), err
}

// CmdRunner describes the signature of the Run function
type CmdRunner func(ctx context.Context, cmd string, options ...RunOption) (int, error)
