package gosh

import (
	"fmt"
	"strings"
)

type Argument string

type ShellCommand interface {
	SubCommands() CommandMap
	Exec([]Argument) error
}

type CommandMap map[string]ShellCommand
type CommandError string

func (c CommandError) Error() string {
	return string(c)
}

func (commands CommandMap) getCompletions(field string) CommandMap {
	completions := make(CommandMap)
	for completion, command := range commands {
		if strings.HasPrefix(completion, field) {
			completions[completion] = command
		}
	}
	return completions
}

func (commands CommandMap) AddCommand(commandName string, command ShellCommand) error {
	if _, ok := commands[commandName]; ok {
		return CommandError(fmt.Sprintf("Command %s is already a top level command", commandName))
	}
	commands[commandName] = command
	return nil
}

func (commands CommandMap) Find(arguments []Argument) (ShellCommand, []Argument, error) {
	var argument Argument
	var i int
	var command ShellCommand

	for i, argument = range arguments {
		nextCommand := commands[string(argument)]
		if nextCommand == nil {
			return nil, nil, CommandError(fmt.Sprintf("No matching command for %v", arguments))
		} else {
			command = nextCommand
			commands = nextCommand.SubCommands()
			if len(commands) == 0 {
				break
			}
		}
	}
	return command, arguments[i+1:], nil
}

func (commands CommandMap) Exec(arguments []Argument) error {
	command, arguments, err := commands.Find(arguments)

	if err != nil {
		return err
	}

	return command.Exec(arguments)
}
