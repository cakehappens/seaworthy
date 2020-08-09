package clioptions

// CommandOptions is part of the functional API for creating new cobra commands
type CommandOptions struct {
	IOStreams
}

// CommandOption is part of the functional API for creating new cobra commands
type CommandOption func(option *CommandOptions)
