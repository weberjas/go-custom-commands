package cec_messages_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCecMessages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CecMessages Suite")
}
