package cli

import (
	"errors"
	"fmt"
	"strings"
)

// App is an abstraction of the options that are chosen
type App struct {
	Name             string
	Description      string
	MaxTerminalWidth int
	margin           int // Internal variable for computing the margin when rendering descriptions
	delimiter        string
	Options          CommandOptions
}

// Run attempts to run the registered command func
func (a *App) Run(cmd string, args ...interface{}) error {
	for _, option := range a.Options {
		if strings.EqualFold(option.Cmd, cmd) {
			return option.Func(args...)
		}
	}
	return fmt.Errorf("Could not find cmd: '%v' for execution", cmd)
}

// Add appends the given command option to the
func (a *App) Add(option ...CommandOption) {
	a.Options.Append(option...)
}

// AddOption allows the user to add an option to the function
func (a *App) AddOption(cmd string, _func OptionFunc) error {
	if cmd == "" {
		return errors.New("no command name was given")
	}
	if _func == nil {
		return errors.New("no function was given")
	}
	a.Options.Append(CommandOption{
		Cmd:  cmd,
		Func: _func,
	})

	return nil
}

// ShowHelp returns a representation of the
func (a *App) ShowHelp() {
	text := ""
	if a.Name != "" {
		stars := padStars(a.MaxTerminalWidth)
		text += stars + "\n\n"

		text += "\n\n" + stars + "\n\n"
	}
	if a.Description != "" {
		text += a.Description + "\n\n"
	}

	if help := a.showCmdHelp(); help != "" {
		text += help
	}

	fmt.Print(text)
}

func (a *App) showCmdHelp() string {
	text := ""
	for _, option := range a.Options {
		desc := strings.Trim(option.Description, " ")
		text += a.reshapeText(desc)
	}
	return text
}

// longestCmdIdx returns the index of the longest command
func (a *App) longestCmdIdx() int {
	max := 0
	for _, cmd := range a.Options {
		if len := len(cmd.Cmd); len > max {
			max = len
		}
	}
	if max > (a.MaxTerminalWidth - 30) {
		max = a.MaxTerminalWidth - 30
	}
	return max
}

// Ensures that the margin is set and "cached" before we proceed with reshaping
// the text
func (a *App) computeMargin() {
	if a.margin == 0 {
		a.margin = a.longestCmdIdx()
		a.delimiter = " - "
	}
}

func (a *App) reshapeText(text string) string {
	a.computeMargin()
	threshold := a.MaxTerminalWidth - a.margin - len(a.delimiter)
	_text := ""
	if len(text) < threshold {
		_text = padSpaces(30) + text
	} else {
		words := strings.Split(text, " ")
		paragraphs := make([]string, 0, 10)
		pt := padSpaces(threshold)
		for _, word := range words {
			if len(pt)+len(word) < threshold {
				pt += word + " "
			} else {
				paragraphs = append(paragraphs, strings.Trim(pt, " "))
				pt = padSpaces(threshold)
			}
		}
		_text = strings.Join(paragraphs, "\n")
	}
	return _text
}

func (a *App) centerPadText(text string) string {
	_len := len(text)
	if _len > a.MaxTerminalWidth {
		return text
	}
	numspaces := (a.MaxTerminalWidth - _len) / 2
	return padSpaces(numspaces) + text + padSpaces(numspaces)
}

func padSpaces(num int) string {
	return padChar(' ', num)
}

func padStars(num int) string {
	return padChar('*', num)
}

func padChar(char rune, num int) string {
	s := ""
	for i := 0; i < num; i++ {
		s += string(char)
	}
	return s
}
