package dashboards_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDashboards(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dashboards Suite")
}
