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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type simpleCommand struct {
	executed  bool
	arguments []string
}

func (t *simpleCommand) Completions() []string {
	return nil
}

func (t *simpleCommand) Exec(arguments []string) error {
	t.executed = true
	t.arguments = arguments
	return nil
}

func newSimpleCommand() *simpleCommand {
	return &simpleCommand{false, nil}
}

/*func newComplexCommand() *complexCommand {
	return &complexCommand{
		map[string]*simpleCommand{
			"subCmd1": newSimpleCommand(),
			"subCmd2": newSimpleCommand(),
		},
		false,
	}
}*/

var _ = Describe("CommandMap", func() {
	Describe("functions", func() {
		var commands CommandMap
		BeforeEach(func() {
			commands = CommandMap{
				"john":  nil,
				"james": nil,
				"mary":  nil,
				"nancy": nil,
			}
		})

		Describe("getCompletions", func() {
			It("should return a CommandMap of all the commands when the field is blank", func() {
				Expect(commands.getCompletions("")).To(Equal(commands))
			})

			It("should return only those commands with matching prefixes", func() {
				Expect(commands.getCompletions("j")).To(Equal(CommandMap{
					"john":  nil,
					"james": nil,
				}))
			})
		})

		Describe("AddCommand", func() {
			It("Should add a new command to the map", func() {
				err := commands.AddCommand("rita", nil)
				Expect(err).To(BeNil())
				Expect(commands.getCompletions("")).To(Equal(CommandMap{
					"john":  nil,
					"james": nil,
					"mary":  nil,
					"nancy": nil,
					"rita":  nil,
				}))
			})

			It("Should return an error instead of adding a duplicate command", func() {
				err := commands.AddCommand("john", nil)
				Expect(err).To(MatchError(ErrDuplicateCommand))
			})
		})
	})

	Describe("Finding a top level command", func() {
		var commands CommandMap
		var cmd *simpleCommand
		BeforeEach(func() {
			cmd = newSimpleCommand()
			commands = CommandMap{
				"cmd": cmd,
			}
		})

		It("should return an error if no command is found", func() {
			_, _, err := commands.Find([]string{"cmd1"})
			Expect(err).To(MatchError(ErrNoMatchingCommand))
		})

		It("should return a matching command", func() {
			execCmd, _, err := commands.Find([]string{"cmd"})
			Expect(err).To(BeNil())
			Expect(execCmd).To(Equal(cmd))
		})

		It("should return an empty argument slice when no arguments are given", func() {
			_, arguments, _ := commands.Find([]string{"cmd"})
			Expect(arguments).To(Equal([]string{}))
		})

		It("should return the arguments to the command when arguments are given", func() {
			_, arguments, _ := commands.Find([]string{"cmd", "arg1", "arg2"})

			Expect(arguments).To(Equal([]string{"arg1", "arg2"}))
		})
	})

	Describe("Finding a sub-command", func() {
		var commands CommandMap
		var tlc TreeCommand
		BeforeEach(func() {
			tlc = NewTreeCommand(CommandMap{
				"subCmd1": newSimpleCommand(),
				"subCmd2": newSimpleCommand(),
			})
			commands = CommandMap{"tlc": tlc}
		})

		It("should return an error for no matching sub-command", func() {
			execCmd, _, err := commands.Find([]string{"tlc", "subCmd3"})
			Expect(err).To(MatchError(ErrNoMatchingCommand))
			Expect(execCmd).To(BeNil())
		})

		It("should return the sub-command", func() {
			execCmd, _, err := commands.Find([]string{"tlc", "subCmd1"})
			Expect(err).To(BeNil())
			Expect(execCmd).To(Equal(tlc.subCommands["subCmd1"]))
		})

		It("should have an empty argument slice for no arguments", func() {
			_, arguments, _ := commands.Find([]string{"tlc", "subCmd1"})
			Expect(arguments).To(Equal([]string{}))
		})

		It("should return the arguments when given", func() {
			_, arguments, _ := commands.Find([]string{"tlc", "subCmd1", "arg1", "arg2"})
			Expect(arguments).To(Equal([]string{"arg1", "arg2"}))
		})
	})

	Describe("Executing a command", func() {
		var commands CommandMap
		var command *simpleCommand

		BeforeEach(func() {
			command = newSimpleCommand()
			commands = CommandMap{"cmd": command}
		})

		It("Should execute the command if found", func() {
			err := commands.Exec([]string{"cmd"})
			Expect(err).To(BeNil())
			Expect(command.executed).To(BeTrue())
		})

		It("Shoud return an error if the command is not found", func() {
			err := commands.Exec([]string{"foo"})
			Expect(err).To(MatchError(ErrNoMatchingCommand))
		})
	})
})
