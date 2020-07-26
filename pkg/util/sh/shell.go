package sh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type RunState struct {
	RunError error
	Ran      bool
	*os.ProcessState
}

type RunOptions struct {
	Args   []string
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

type RunOption func(o *RunOptions)

// Run provides a functional API for running a process
// Default Behaviors:
// stdout, stderr, stdin are sent to /dev/null
// all environment variables are passed to subprocess
// additionally, you may add vars to RunOptions.Env
//
// RunState is passed to the channel as soon as cmd.Start is called
// StartErr is the error returned from cmd.Start() https://golang.org/pkg/os/exec/#Cmd.Start
// WaitErr is the error from cmd.Wait() https://golang.org/pkg/os/exec/#Cmd.Wait
func Run(ctx context.Context, cmd string, options ...RunOption) *RunState {
	runOpts := &RunOptions{}

	runOpts.Env = make(map[string]string)

	for _, item := range os.Environ() {
		split := strings.SplitN(item, "=", 2)
		runOpts.Env[split[0]] = split[1]
	}

	for _, optFn := range options {
		optFn(runOpts)
	}

	if runOpts.Args == nil {
		runOpts.Args = make([]string, 0)
	}

	cmdEnv := make([]string, len(runOpts.Env))

	for k, v := range runOpts.Env {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v))
	}

	c := exec.CommandContext(ctx, cmd, runOpts.Args...)
	c.Env = cmdEnv

	//fmt.Printf("cmd: %s %s\n", cmd, strings.Join(runOpts.Args, " "))

	c.Stdout = runOpts.Stdout
	c.Stderr = runOpts.Stderr
	c.Stdin = runOpts.Stdin

	runState := &RunState{}

	runState.RunError = c.Run()
	runState.ProcessState = c.ProcessState

	if runState.RunError == nil {
		runState.Ran = true
	} else if _, ok := runState.RunError.(*exec.ExitError); ok {
		runState.Ran = true
	}

	return runState
}
