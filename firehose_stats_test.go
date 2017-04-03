package main_test

import (
	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/cloudfoundry/firehose-stats"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//	"github.com/cloudfoundry/firehose-plugin/testhelpers"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"

	"strings"
)

var _ = Describe("FirehoseStatsPlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *fakes.FakeCliConnection
		var firehoseStatsCmd *FirehoseStatsCmd
		//		var fakeFirehose *testhelpers.FakeFirehose

		PIt("displays debug info when debug flag is passed", func() {

			outputChan := make(chan []string)
			go func() {
				output := io_helpers.CaptureOutput(func() {
					firehoseStatsCmd.Run(fakeCliConnection, []string{"firehose-stats", "--debug"})
				})
				outputChan <- output
			}()

			var output []string
			Eventually(outputChan, 2).Should(Receive(&output))
			outputString := strings.Join(output, " | ")

			Expect(outputString).To(ContainSubstring("Hit Ctrl+c to exit"))
			Expect(outputString).To(ContainSubstring("WEBSOCKET REQUEST"))
			Expect(outputString).To(ContainSubstring("WEBSOCKET RESPONSE"))
		})
	})
})
