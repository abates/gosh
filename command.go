/**
 * Copyright 2015 Andrew Bates
 *
 * Licensed under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package gosh

import (
	"os"
	"strings"
)

// Completable is the interface for making a Command auto-completable
//
// Completions returns a slice of strings representing the list of completion
// candidates.  The method is called with the field immediately following the
// command.  For instance.  If the command being completed is "ls
// /has/cheesburger" Then the Completion method for ls will be provided the
// string /has/cheeseburger
type Completable interface {
	Completions(field string) []string
}

// TreeCommand is a concrete implementation of Command
//
// TreeCommand provides the ability to create a hierarchy of commands.  This
// type of command hierarchy is very common in command line interfaces for
// network appliances such as router and firewalls (think JunOS or Cisco IOS)
type TreeCommand struct {
	subCommands CommandMap
}

// SubCommands returns the CommandMap of sub commands that belong to this
// TreeCommand
func (t TreeCommand) SubCommands() CommandMap {
	return t.subCommands
}

// Exec does nothing since a TreeCommand only contains sub-commands
func (t TreeCommand) Exec() error {
	return nil
}

// Add another sub-command to this TreeCommand
func (t TreeCommand) Add(name string, command Command) error {
	return t.subCommands.Add(name, command)
}

// NewTreeCommand creates a TreeCommand for the given CommandMap
func NewTreeCommand(commands CommandMap) TreeCommand {
	tree := TreeCommand{
		subCommands: commands,
	}
	return tree
}

// Command indicates that an object can be executed
//
// Exec should perform any computation necessary to execute the command that
// provides the interface.  Prior to calling the Exec method, the shell will assign
// the arguments to os.Args.  If the command is expecting any arguments, then they will
// be availabe as os.Args.  The argument list includes the command path as os.Args[0]
type Command interface {
	Exec() error
}

// CommandMap is exactly what it sounds like.
//
// CommandMap is a map of Commands that are keyed by the command name that
// should be typed at the prompt.  If the prompt should execute a command when
// ls\n is typed, then the command name is ls and the Command is the concrete
// implementation that provides Exec.  When ls\n is typed at the prompt, the
// backing Command will be looked up in the CommandMap and its Exec method
// called
type CommandMap map[string]Command

func (commands CommandMap) getCompletions(field string) CommandMap {
	completions := make(CommandMap)
	for completion, command := range commands {
		if strings.HasPrefix(completion, field) {
			completions[completion] = command
		}
	}
	return completions
}

// Add a comand to the map
func (commands CommandMap) Add(commandName string, command Command) error {
	if _, ok := commands[commandName]; ok {
		return ErrDuplicateCommand
	}
	commands[commandName] = command
	return nil
}

// Find traverses the command map using the arguments slice and return the
// Command whose path exactly matches the argument list.  If no Command can be
// found with an exact matching path then ErrNoMatchingCommand is returned.
func (commands CommandMap) Find(arguments []string) (Command, []string, error) {
	var argument string
	var i int
	var command Command

	for i, argument = range arguments {
		nextCommand := commands[argument]
		if nextCommand == nil {
			return nil, nil, ErrNoMatchingCommand
		}

		command = nextCommand
		if nextCommand, ok := nextCommand.(TreeCommand); ok {
			commands = nextCommand.SubCommands()
			if len(commands) == 0 {
				break
			}
		} else {
			break
		}
	}
	return command, arguments[i+1:], nil
}

// Exec finds and execute a command corresponding to the argument list
func (commands CommandMap) Exec(fields []string) error {
	command, arguments, err := commands.Find(fields)

	if err != nil {
		return err
	}

	oldArgs := make([]string, len(os.Args))
	copy(oldArgs, os.Args)

	defer func() { os.Args = oldArgs }()
	os.Args = make([]string, len(arguments)+1)
	os.Args[0] = strings.Join(fields[:len(fields)-len(arguments)], " ")
	for i, argument := range arguments {
		os.Args[i+1] = argument
	}
	return command.Exec()
}
