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

var _ = Describe("completer", func() {
	Describe("Single completion", func() {
		var c *completer
		BeforeEach(func() {
			c = newCompleter(CommandMap{
				"cmd": newTestCommand(),
			})
		})

		It("Should auto complete an empty string with the whole command", func() {
			head, completions, tail := c.complete("", 0)
			Expect(head).To(Equal(""))
			Expect(completions).To(Equal([]string{"cmd"}))
			Expect(tail).To(Equal(""))
		})

		It("Should auto complete a given string with the whole command", func() {
			head, completions, tail := c.complete("c", 1)
			Expect(head).To(Equal(""))
			Expect(completions).To(Equal([]string{"cmd"}))
			Expect(tail).To(Equal(""))
		})
	})

	Describe("top level behavior", func() {
		var c *completer
		BeforeEach(func() {
			c = newCompleter(CommandMap{
				"john":  newTestCommand(),
				"james": newTestCommand(),
				"mary":  newTestCommand(),
				"nancy": newTestCommand(),
			})
		})

		It("Should return all the top level strings when the empty string is supplied", func() {
			wanted := []string{"james", "john", "mary", "nancy"}
			head, completions, tail := c.complete("", 0)
			Expect(head).To(Equal(""))
			Expect(completions).To(Equal(wanted))
			Expect(tail).To(Equal(""))
		})

		It("Should return strings that match the input prefix", func() {
			wanted := []string{"james", "john"}
			head, completions, tail := c.complete("j", 1)
			Expect(head).To(Equal(""))
			Expect(completions).To(Equal(wanted))
			Expect(tail).To(Equal(""))
		})
	})

	Describe("Second level HierarchyCompleter response", func() {
		var c *completer
		var commands CommandMap
		BeforeEach(func() {
			commands = CommandMap{
				"johnson": newTestCommand(),
				"john": NewTreeCommand(CommandMap{
					"jacob":        newTestCommand(),
					"jingleheimer": newTestCommand(),
					"schmidt":      newTestCommand(),
				}),
				"james": newTestCommand(),
				"mary":  newTestCommand(),
				"nancy": newTestCommand(),
			}

			c = newCompleter(commands)
		})

		It("Should return all the second level tokens when there is an exact match for the first field and no second field", func() {
			wanted := []string{"jacob", "jingleheimer", "schmidt"}
			head, completions, tail := c.complete("john ", 5)
			Expect(head).To(Equal("john "))
			Expect(completions).To(Equal(wanted))
			Expect(tail).To(Equal(""))
		})

		It("Should return only matching second level tokens when there is an exact match for the first field and second field", func() {
			wanted := []string{"jacob", "jingleheimer"}
			head, completions, tail := c.complete("john j", 6)
			Expect(head).To(Equal("john "))
			Expect(completions).To(Equal(wanted))
			Expect(tail).To(Equal(""))
		})

		It("Should not return parent completions when there is ambiguity", func() {
			wanted := []string{"jacob", "jingleheimer"}
			head, completions, tail := c.complete("john j", 6)
			Expect(head).To(Equal("john "))
			Expect(completions).To(Equal(wanted))
			Expect(tail).To(Equal(""))
		})
	})

	Describe("Simple command completions", func() {
		var command *testCommand
		var c *completer
		var completions []string

		BeforeEach(func() {
			command = newTestCommand()
			completions = []string{
				"aarg1",
				"aarg2",
				"barg1",
				"barg2",
			}

			command.setCompletions(completions)
			c = newCompleter(CommandMap{
				"cmd": command,
			})
		})

		It("should return all arguments when completing the command with no prefix", func() {
			head, completions, tail := c.complete("cmd ", 4)
			Expect(head).To(Equal("cmd "))
			Expect(completions).To(Equal([]string{
				"aarg1",
				"aarg2",
				"barg1",
				"barg2",
			}))
			Expect(tail).To(Equal(""))
		})

		It("should return matching arguments for a given prefix", func() {
			head, completions, tail := c.complete("cmd a", 5)
			Expect(head).To(Equal("cmd "))
			Expect(completions).To(Equal([]string{
				"aarg1",
				"aarg2",
			}))
			Expect(tail).To(Equal(""))
		})
	})
})
