package clioptions

import (

)

type CommandOptions struct {
	IOStreams
}

type CommandOption func(option *CommandOptions)
