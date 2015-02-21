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

var _ = Describe("Completer", func() {
	Describe("top level behavior", func() {
		var completer *Completer
		BeforeEach(func() {
			completer = NewCompleter(CommandMap{
				"john":  newTestCommand(),
				"james": newTestCommand(),
				"mary":  newTestCommand(),
				"nancy": newTestCommand(),
			})
		})

		It("Should return all the top level strings when the empty string is supplied", func() {
			wanted := []string{"james", "john", "mary", "nancy"}
			Expect(completer.Complete("")).To(Equal(wanted))
		})

		It("Should return strings that match the input prefix", func() {
			wanted := []string{"james", "john"}
			Expect(completer.Complete("j")).To(Equal(wanted))
		})
	})

	Describe("Second level HierarchyCompleter response", func() {
		var completer *Completer
		BeforeEach(func() {
			completer = NewCompleter(CommandMap{
				"john": NewTreeCommand(CommandMap{
					"jacob":        newTestCommand(),
					"jingleheimer": newTestCommand(),
					"schmidt":      newTestCommand(),
				}),
				"james": newTestCommand(),
				"mary":  newTestCommand(),
				"nancy": newTestCommand(),
			})
		})

		It("Should return all the second level tokens when there is an exact match for the first field and no second field", func() {
			wanted := []string{"john jacob", "john jingleheimer", "john schmidt"}
			Expect(completer.Complete("john ")).To(Equal(wanted))
		})

		It("Should return only matching second level tokens when there is an exact match for the first field and second field", func() {
			wanted := []string{"john jacob", "john jingleheimer"}
			Expect(completer.Complete("john j")).To(Equal(wanted))
		})
	})

	Describe("Simple command completions", func() {
		var command *testCommand
		var completer *Completer
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
			completer = NewCompleter(CommandMap{
				"cmd": command,
			})
		})

		It("should return all arguments when completing the command with no prefix", func() {
			Expect(completer.Complete("cmd ")).To(Equal([]string{
				"cmd aarg1",
				"cmd aarg2",
				"cmd barg1",
				"cmd barg2",
			}))
		})

		It("should return matching arguments for a given prefix", func() {
			Expect(completer.Complete("cmd a")).To(Equal([]string{
				"cmd aarg1",
				"cmd aarg2",
			}))
		})
	})
})
