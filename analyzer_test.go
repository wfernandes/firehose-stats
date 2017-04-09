package main_test

import (
	"testing"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/gomega"
	. "github.com/wfernandes/firehose-stats"
)

var (
	envelopes chan *events.Envelope
	printer   *mockPrinter
	analyzer  *Analyzer
)

func setup(t *testing.T) {
	RegisterTestingT(t)
	envelopes = make(chan *events.Envelope, 100)
	printer = &mockPrinter{}
	analyzer = NewAnalyzer(
		envelopes,
		printer,
		WithPrintInterval(10*time.Millisecond),
	)
}

func TestAnalyzer_CountsTotal(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 5; i++ {
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	f := func() int {
		return printer.Total
	}
	Eventually(f).Should(Equal(5))
}

func TestAnalyzer_CountsLogMessages(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 5; i++ {
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	f := func() int {
		return printer.LogMessages
	}
	Eventually(f).Should(Equal(5))
}

func TestAnalyzer_CountsCounterEvents(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 3; i++ {
		envelopes <- buildEnvelope(events.Envelope_CounterEvent)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	getCounterEvents := func() int {
		return printer.CounterEvents
	}
	getTotal := func() int {
		return printer.Total
	}
	Eventually(getCounterEvents).Should(Equal(3))
	Eventually(getTotal).Should(Equal(6))
}

func TestAnalyzer_CountsContainerMetrics(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 3; i++ {
		envelopes <- buildEnvelope(events.Envelope_ContainerMetric)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	getContainerMetrics := func() int {
		return printer.ContainerMetrics
	}
	getTotal := func() int {
		return printer.Total
	}
	Eventually(getContainerMetrics).Should(Equal(3))
	Eventually(getTotal).Should(Equal(6))
}

func TestAnalyzer_CountsValueMetrics(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 4; i++ {
		envelopes <- buildEnvelope(events.Envelope_ValueMetric)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	getValueMetrics := func() int {
		return printer.ValueMetrics
	}
	getTotal := func() int {
		return printer.Total
	}
	Eventually(getValueMetrics).Should(Equal(4))
	Eventually(getTotal).Should(Equal(8))
}

func TestAnalyzer_AvgEnvelopeSize(t *testing.T) {
	setup(t)

	go analyzer.Start()
	e1 := &events.Envelope{
		Origin:    proto.String("some-origin"),
		EventType: events.Envelope_LogMessage.Enum(),
		LogMessage: &events.LogMessage{
			Message:     []byte("some-random-msg"),
			MessageType: events.LogMessage_OUT.Enum(),
			Timestamp:   proto.Int64(time.Now().UnixNano()),
		},
	}
	e2 := &events.Envelope{
		Origin:    proto.String("some-other-origin"),
		EventType: events.Envelope_LogMessage.Enum(),
		LogMessage: &events.LogMessage{
			Message:     []byte("another-random-msg"),
			MessageType: events.LogMessage_OUT.Enum(),
			Timestamp:   proto.Int64(time.Now().UnixNano()),
		},
	}
	for i := 0; i < 5; i++ {
		envelopes <- e1
	}
	envelopes <- e2
	getAvgEnvelopeSize := func() int {
		return printer.AvgEnvelopeSize
	}
	getTotal := func() int {
		return printer.Total
	}
	expectedAvg := ((e1.Size() * 5) + e2.Size()) / 6
	Eventually(getTotal).Should(Equal(6))
	Eventually(getAvgEnvelopeSize).Should(Equal(expectedAvg))
}

func buildEnvelope(envtype events.Envelope_EventType) *events.Envelope {
	return &events.Envelope{
		Origin:    proto.String("some-origin"),
		EventType: &envtype,
	}
}

type mockPrinter struct {
	Total            int
	LogMessages      int
	CounterEvents    int
	ContainerMetrics int
	ValueMetrics     int
	AvgEnvelopeSize  int
}

func newMockPrinter() *mockPrinter {
	return &mockPrinter{}
}

func (p *mockPrinter) Print(s Stats) {
	p.Total = int(s["TotalMessages"])
	p.LogMessages = int(s["LogMessages"])
	p.CounterEvents = int(s["CounterEvents"])
	p.ContainerMetrics = int(s["ContainerMetrics"])
	p.ValueMetrics = int(s["ValueMetrics"])
	p.AvgEnvelopeSize = int(s["AvgEnvelopeSize"])
}
