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
	"fmt"
	"strings"
)

type ShellCommand interface {
	SubCommands() CommandMap
	Exec([]string) error
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

func (commands CommandMap) Find(arguments []string) (ShellCommand, []string, error) {
	var argument string
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

func (commands CommandMap) Exec(arguments []string) error {
	command, arguments, err := commands.Find(arguments)

	if err != nil {
		return err
	}

	return command.Exec(arguments)
}
