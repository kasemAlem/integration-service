package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNudgeMigrate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NudgeMigrate Test Suite")
}
