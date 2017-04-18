package main_test

import (
	"sync/atomic"
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
	printer = &mockPrinter{
		Counts: make(chan Stats, 100),
	}
	analyzer = NewAnalyzer(
		envelopes,
		printer,
		WithPrintInterval(10*time.Millisecond),
	)
}

func TestCountsTotal(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 500; i++ {
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	f := func() int64 {
		return atomic.LoadInt64(&printer.Total)
	}
	Eventually(f).Should(BeEquivalentTo(500))
}

func TestCountsLogMessages(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 5; i++ {
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	f := func() int64 {
		return atomic.LoadInt64(&printer.LogMessages)
	}
	Eventually(f).Should(BeEquivalentTo(5))
}

func TestCountsCounterEvents(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 3; i++ {
		envelopes <- buildEnvelope(events.Envelope_CounterEvent)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	getCounterEvents := func() int64 {
		return atomic.LoadInt64(&printer.CounterEvents)
	}
	getTotal := func() int64 {
		return atomic.LoadInt64(&printer.Total)
	}
	Eventually(getCounterEvents).Should(BeEquivalentTo(3))
	Eventually(getTotal).Should(BeEquivalentTo(6))
}

func TestCountsContainerMetrics(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 3; i++ {
		envelopes <- buildEnvelope(events.Envelope_ContainerMetric)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	getContainerMetrics := func() int64 {
		return atomic.LoadInt64(&printer.ContainerMetrics)
	}
	getTotal := func() int64 {
		return atomic.LoadInt64(&printer.Total)
	}
	Eventually(getContainerMetrics).Should(BeEquivalentTo(3))
	Eventually(getTotal).Should(BeEquivalentTo(6))
}

func TestCountsValueMetrics(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 4; i++ {
		envelopes <- buildEnvelope(events.Envelope_ValueMetric)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
	}

	getValueMetrics := func() int64 {
		return atomic.LoadInt64(&printer.ValueMetrics)
	}
	getTotal := func() int64 {
		return atomic.LoadInt64(&printer.Total)
	}
	Eventually(getValueMetrics).Should(BeEquivalentTo(4))
	Eventually(getTotal).Should(BeEquivalentTo(8))
}

func TestTotalEnvelopeSize(t *testing.T) {
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
	getTotalEnvelopeSize := func() int64 {
		return atomic.LoadInt64(&printer.TotalEnvelopeSize)
	}
	getTotal := func() int64 {
		return atomic.LoadInt64(&printer.Total)
	}
	expectedTotal := (e1.Size() * 5) + e2.Size()
	Eventually(getTotal).Should(BeEquivalentTo(6))
	Eventually(getTotalEnvelopeSize).Should(BeEquivalentTo(expectedTotal))
}

func TestTotalCountEqualsSumOfEnvelopes(t *testing.T) {
	setup(t)

	go analyzer.Start()
	for i := 0; i < 4; i++ {
		envelopes <- buildEnvelope(events.Envelope_ValueMetric)
		envelopes <- buildEnvelope(events.Envelope_LogMessage)
		envelopes <- buildEnvelope(events.Envelope_CounterEvent)
		envelopes <- buildEnvelope(events.Envelope_ContainerMetric)
		envelopes <- buildEnvelope(events.Envelope_HttpStartStop)
	}

	var s Stats
	Eventually(printer.Counts, "2s").Should(Receive(&s))
	Expect(s["TotalMessages"]).To(
		Equal(
			s["LogMessages"] + s["CounterEvents"] + s["ContainerMetrics"] + s["ValueMetrics"] + s["HttpStartStops"],
		),
	)
}

func buildEnvelope(envtype events.Envelope_EventType) *events.Envelope {
	return &events.Envelope{
		Origin:    proto.String("some-origin"),
		EventType: &envtype,
	}
}

type mockPrinter struct {
	Total             int64
	LogMessages       int64
	CounterEvents     int64
	ContainerMetrics  int64
	ValueMetrics      int64
	HttpStartStops    int64
	TotalEnvelopeSize int64
	Counts            chan Stats
}

func (p *mockPrinter) Print(s Stats) {
	atomic.StoreInt64(&p.Total, s["TotalMessages"])
	atomic.StoreInt64(&p.LogMessages, s["LogMessages"])
	atomic.StoreInt64(&p.CounterEvents, s["CounterEvents"])
	atomic.StoreInt64(&p.ContainerMetrics, s["ContainerMetrics"])
	atomic.StoreInt64(&p.ValueMetrics, s["ValueMetrics"])
	atomic.StoreInt64(&p.HttpStartStops, s["HttpStartStops"])
	atomic.StoreInt64(&p.TotalEnvelopeSize, s["TotalEnvelopeSize"])
	p.Counts <- s
}
