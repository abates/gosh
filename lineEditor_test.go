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
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

type testResponse struct {
	resp string
	err  error
}

type testLineEditor struct {
	response    chan testResponse
	closeErr    error
	closeCalled bool
}

func (t *testLineEditor) Prompt(string) (string, error) {
	var resp testResponse
	select {
	case resp = <-t.response:
	default:
		resp = testResponse{"", nil}
	}
	return resp.resp, resp.err
}

func (t *testLineEditor) Close() error {
	t.closeCalled = true
	return t.closeErr
}

func (t *testLineEditor) addResponse(resp string, err error) {
	t.response <- testResponse{resp: resp, err: err}
}

func (t *testLineEditor) end() {
	t.response <- testResponse{resp: "", err: io.EOF}
}

func newTestLineEditor() *testLineEditor {
	return &testLineEditor{
		response:    make(chan testResponse, 20),
		closeErr:    nil,
		closeCalled: false,
	}
}

type nonCloseableLineEditor struct{}

func (t *nonCloseableLineEditor) Prompt(string) (string, error) { return "", nil }

var _ = Describe("DefaultLineEditor", func() {
	It("Should append the last command to the history", func() {
		var b bytes.Buffer

		stdin_r, stdin_wr, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = stdin_r

		editor := NewDefaultLineEditor(CommandMap{})

		stdin_wr.Write([]byte("cmd\n"))
		editor.Prompt(">")
		os.Stdin = oldStdin
		editor.liner.WriteHistory(&b)
		Expect(b.String()).To(Equal("cmd\n"))
	})
})
