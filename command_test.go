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
	"os"
)

type testCommand struct {
	completions  []string
	executed     bool
	arguments    []string
	execErr      error
	execCallback func() error
}

func (t *testCommand) Completions(substring string) []string {
	return t.completions
}

func (t *testCommand) Exec() error {
	t.executed = true
	t.arguments = make([]string, len(os.Args))
	copy(t.arguments, os.Args)
	if t.execCallback != nil {
		return t.execCallback()
	}
	return t.execErr
}

func (t *testCommand) setCompletions(completions []string) {
	t.completions = completions
}

func newTestCommand() *testCommand {
	return &testCommand{
		completions:  nil,
		executed:     false,
		arguments:    nil,
		execErr:      nil,
		execCallback: nil,
	}
}

func newCallbackCommand(callback func() error) *testCommand {
	return &testCommand{
		completions:  nil,
		executed:     false,
		arguments:    nil,
		execErr:      nil,
		execCallback: callback,
	}
}

var _ = Describe("CommandMap", func() {
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

	Describe("Add", func() {
		It("Should add a new command to the map", func() {
			Expect(commands.Add("rita", nil)).To(Succeed())
			Expect(commands.getCompletions("")).To(Equal(CommandMap{
				"john":  nil,
				"james": nil,
				"mary":  nil,
				"nancy": nil,
				"rita":  nil,
			}))
		})

		It("Should return an error instead of adding a duplicate command", func() {
			err := commands.Add("john", nil)
			Expect(err).To(MatchError(ErrDuplicateCommand))
		})
	})

	Describe("Finding a top level command", func() {
		var cmd *testCommand
		BeforeEach(func() {
			cmd = newTestCommand()
			commands.Add("cmd", cmd)
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

	Describe("Exec", func() {
		It("Should return an error if executing a command that can't be found", func() {
			Expect(commands.Exec([]string{"invalid"})).To(MatchError(ErrNoMatchingCommand))
		})

		It("Should set os.Args to the program name and the argument list", func() {
			cmd := newCallbackCommand(func() error {
				Expect(os.Args).To(Equal([]string{"cmd", "arg1", "arg2"}))
				return nil
			})
			commands.Add("cmd", cmd)
			oldArgs := os.Args
			Expect(commands.Exec([]string{"cmd", "arg1", "arg2"})).To(Succeed())
			Expect(os.Args).To(Equal(oldArgs))
		})
	})
})

var _ = Describe("TreeCommand", func() {
	Describe("Finding a sub-command", func() {
		var commands CommandMap
		var tlc TreeCommand
		BeforeEach(func() {
			tlc = NewTreeCommand(CommandMap{
				"subCmd1": newTestCommand(),
				"subCmd2": newTestCommand(),
			})
			commands = CommandMap{"tlc": tlc}
		})

		It("should return the top level command if it is a TreeCommand with no sub commands", func() {
			treeCmd := NewTreeCommand(CommandMap{})
			commands.Add("treeCmd", treeCmd)
			cmd, arguments, _ := commands.Find([]string{"treeCmd", "arg", "arg2"})
			Expect(cmd).To(Equal(treeCmd))
			Expect(arguments).To(Equal([]string{"arg", "arg2"}))
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

		It("Should return nil when executing", func() {
			Expect(tlc.Exec()).To(BeNil())
		})

		It("Should set os.Args[0] to the full command path when executing the command", func() {
			tlc.Add("subCmd3", newCallbackCommand(func() error {
				Expect(os.Args).To(Equal([]string{"tlc subCmd3", "arg1", "arg2"}))
				return nil
			}))
			Expect(commands.Exec([]string{"tlc", "subCmd3", "arg1", "arg2"})).To(Succeed())
		})
	})
})
