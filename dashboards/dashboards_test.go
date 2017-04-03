package dashboards_test

import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	"github.com/wfernandes/firehose-stats/dashboards"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dashboards", func() {

	Context("General", func() {

		var generalDash dashboards.Dashboard

		BeforeEach(func() {
			generalDash := dashboards.NewGeneral()
		})

		It("creates charts upon initialization", func() {
			Expect(generalDash.GetCharts()).To(HaveLen(3))
		})

		It("sifts events for its charts", func() {
			envelope := events.Envelope{
				Origin:     proto.String("origin"),
				Timestamp:  proto.Int64(1000000000),
				Deployment: proto.String("deployment-name"),
				Job:        proto.String("doppler"),
				EventType:  events.Envelope_LogMessage.Enum(),
				LogMessage: &events.LogMessage{
					Message:     []byte("some log message"),
					MessageType: events.LogMessage_OUT.Enum(),
					Timestamp:   proto.Int64(1000000000),
				},
			}

			generalDash.Consume(&envelope)

		})

	})
})
