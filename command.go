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
	"strings"
)

type TreeCommand struct {
	subCommands CommandMap
	completions []string
}

func (t TreeCommand) Completions() []string {
	return t.completions
}

func (t TreeCommand) SubCommands() CommandMap {
	return t.subCommands
}

func (t TreeCommand) Exec([]string) error {
	return nil
}

func NewTreeCommand(commands CommandMap) TreeCommand {
	tree := TreeCommand{
		subCommands: commands,
		completions: make([]string, len(commands)),
	}

	i := 0
	for commandName, _ := range commands {
		tree.completions[i] = commandName
		i += 1
	}
	return tree
}

type Command interface {
	Exec([]string) error
	Completions() []string
}

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

func (commands CommandMap) Add(commandName string, command Command) error {
	if _, ok := commands[commandName]; ok {
		return ErrDuplicateCommand
	}
	commands[commandName] = command
	return nil
}

func (commands CommandMap) Find(arguments []string) (Command, []string, error) {
	var argument string
	var i int
	var command Command

	for i, argument = range arguments {
		nextCommand := commands[argument]
		if nextCommand == nil {
			return nil, nil, ErrNoMatchingCommand
		} else {
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
	}
	return command, arguments[i+1:], nil
}

func (commands CommandMap) Exec(arguments []string) error {
	command, arguments, err := commands.Find(arguments)

	if err != nil {
		return err
	}

	return command.Exec(arguments)
}
