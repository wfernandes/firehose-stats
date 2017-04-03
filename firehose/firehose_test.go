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
	var consumer *mockConsumer

	BeforeEach(func() {
		lineCounter = 0
		printer = new(fakes.FakePrinter)
		stdout = ""
		consumer = newMockConsumer()
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			stdout += fmt.Sprintf(format, a...)
			lineCounter++
			return len(stdout), nil
		}
		stdin = &fakeStdin{[]byte{'\n'}, false}
		ui = terminal.NewUI(stdin, printer)
	})

	Context("Start", func() {
		Context("when the connection to doppler cannot be established", func() {
			It("shows a meaningful error", func() {
				client := firehose.NewClient("invalidToken", "badEndpoint", consumer, ui)
				client.Start()
				Expect(stdout).To(ContainSubstring("Error dialing traffic controller server"))
			})
		})

		Context("when the connection to doppler works", func() {
			var fakeFirehose *testhelpers.FakeFirehose
			BeforeEach(func() {
				fakeFirehose = testhelpers.NewFakeFirehose("ACCESS_TOKEN")
				fakeFirehose.SendEvent(events.Envelope_LogMessage, "This is a very special test message")
				fakeFirehose.SendEvent(events.Envelope_ValueMetric, "This is a valuemetric")
				fakeFirehose.SendEvent(events.Envelope_CounterEvent, "This is a counterevent")
				fakeFirehose.SendEvent(events.Envelope_ContainerMetric, "This is a containermetric")
				fakeFirehose.Start()
			})

			It("processes events", func() {
				client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), consumer, ui)
				client.Start()
				Eventually(consumer.ConsumeCalled).Should(Receive(BeTrue()))
				var receivedEnvelope *events.Envelope
				Eventually(consumer.ConsumeInput.e).Should(Receive(&receivedEnvelope))
				Expect(receivedEnvelope.String()).To(ContainSubstring("This is a very special test message"))

				Eventually(consumer.ConsumeInput.e).Should(Receive(&receivedEnvelope))
				Expect(receivedEnvelope.String()).To(ContainSubstring("This is a valuemetric"))

				Eventually(consumer.ConsumeInput.e).Should(Receive(&receivedEnvelope))
				Expect(receivedEnvelope.String()).To(ContainSubstring("This is a counterevent"))

				Eventually(consumer.ConsumeInput.e).Should(Receive(&receivedEnvelope))
				Expect(receivedEnvelope.String()).To(ContainSubstring("This is a containermetric"))
			})
		})
	})
})
