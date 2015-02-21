package gosh

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	"io"
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

var _ = Describe("LineEditor", func() {

})
