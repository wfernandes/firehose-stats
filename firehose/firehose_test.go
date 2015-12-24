package firehose_test

import (
	"fmt"
	"io"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/terminal/fakes"
	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wfernandes/firehose-stats/firehose"
	"github.com/wfernandes/firehose-stats/testhelpers"
)

type fakeStdin struct {
	Input []byte
	done  bool
}

func (r *fakeStdin) Read(p []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}
	for i, b := range r.Input {
		p[i] = b
	}
	r.done = true
	return len(r.Input), nil
}

var _ = Describe("Firehose", func() {
	var printer *fakes.FakePrinter
	var ui terminal.UI
	var stdin *fakeStdin
	var stdout string
	var lineCounter int

	BeforeEach(func() {
		lineCounter = 0
		printer = new(fakes.FakePrinter)
		stdout = ""
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			stdout += fmt.Sprintf(format, a...)
			lineCounter++
			return len(stdout), nil
		}
		stdin = &fakeStdin{[]byte{'\n'}, false}
		ui = terminal.NewUI(stdin, printer)
	})

	Context("Start", func() {
		FContext("when the connection to doppler cannot be established", func() {
			It("shows a meaningful error", func() {
				client := firehose.NewClient("invalidToken", "badEndpoint", ui)
				client.Start()
				Eventually(stdout).Should(ContainSubstring("Error dialing traffic controller server"))
			})

		})

		Context("when the connection to doppler works", func() {
			var fakeFirehose *testhelpers.FakeFirehose
			BeforeEach(func() {
				fakeFirehose = testhelpers.NewFakeFirehose("ACCESS_TOKEN")
				fakeFirehose.SendEvent(events.Envelope_LogMessage, "This is a very special test message")
				fakeFirehose.SendEvent(events.Envelope_ValueMetric, "valuemetric")
				fakeFirehose.SendEvent(events.Envelope_CounterEvent, "counterevent")
				fakeFirehose.SendEvent(events.Envelope_ContainerMetric, "containermetric")
				fakeFirehose.SendEvent(events.Envelope_Error, "this is an error")
				fakeFirehose.SendEvent(events.Envelope_HttpStart, "start request")
				fakeFirehose.SendEvent(events.Envelope_HttpStop, "stop request")
				fakeFirehose.SendEvent(events.Envelope_HttpStartStop, "startstop request")
				fakeFirehose.Start()
			})
			It("prints out debug information if demanded", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("WEBSOCKET REQUEST"))
				Expect(stdout).To(ContainSubstring("WEBSOCKET RESPONSE"))
			})

			It("prints out log messages to the terminal", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("This is a very special test message"))
			})

		})
	})
})
