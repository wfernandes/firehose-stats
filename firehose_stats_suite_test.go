package stats_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFirehoseStats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FirehoseStats Suite")
}
