package gosh

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGosh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gosh Suite")
}
