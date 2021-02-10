package cli

// CommandOption is an option that is chosen
type CommandOption struct {
	Cmd         string // Cmd is the command that the given option is associated with
	Description string
	Func        OptionFunc
}

// CommandOptions is a convenient alias for a collection of command options
type CommandOptions []CommandOption

// AsMap returns the commandoption represented as a map of string indices -
// identified based on the command
func (co *CommandOptions) AsMap() map[string]CommandOption {
	_map := make(map[string]CommandOption)
	for _, option := range *co {
		_map[option.Cmd] = option
	}
	return _map
}

// Append is a helper that makes appending functions easy
func (co *CommandOptions) Append(opts ...CommandOption) {
	if co == nil || len(*co) == 0 {
		*co = make(CommandOptions, 0, 5)
	}
	*co = append(*co, opts...)
}

// OptionFunc is a function template, of sorts that helps clients of this code
// register functions to perform the
type OptionFunc func(args ...interface{}) error
