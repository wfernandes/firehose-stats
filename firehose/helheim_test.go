package firehose_test

import "github.com/cloudfoundry/sonde-go/events"

type mockConsumer struct {
	ConsumeCalled chan bool
	ConsumeInput  struct {
		e chan *events.Envelope
	}
}

func newMockConsumer() *mockConsumer {
	m := &mockConsumer{}
	m.ConsumeCalled = make(chan bool, 100)
	m.ConsumeInput.e = make(chan *events.Envelope, 100)
	return m
}
func (m *mockConsumer) Consume(e *events.Envelope) {
	m.ConsumeCalled <- true
	m.ConsumeInput.e <- e
}
