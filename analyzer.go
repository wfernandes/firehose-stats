package main

import (
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

type Printer interface {
	Print(Stats)
}

type Stats map[string]int64

type Analyzer struct {
	messages      <-chan *events.Envelope
	printer       Printer
	printInterval time.Duration
	stats         Stats
	mu            sync.RWMutex
}

func NewAnalyzer(msg <-chan *events.Envelope, p Printer, opts ...AnalyzerOpts) *Analyzer {
	a := &Analyzer{
		messages:      msg,
		printer:       p,
		printInterval: 5 * time.Second,
		stats:         Stats{},
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

func (a *Analyzer) Start() {
	ticker := time.NewTicker(a.printInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			a.mu.RLock()
			data := Stats{}
			for k, v := range a.stats {
				data[k] = v
			}
			a.mu.RUnlock()
			a.printer.Print(data)
		}
	}()

	for e := range a.messages {
		a.add("TotalEnvelopeSize", int64(e.Size()))
		switch e.GetEventType() {
		case events.Envelope_CounterEvent:
			a.incr("CounterEvents")
		case events.Envelope_LogMessage:
			a.incr("LogMessages")
		case events.Envelope_ContainerMetric:
			a.incr("ContainerMetrics")
		case events.Envelope_ValueMetric:
			a.incr("ValueMetrics")
		}
	}
}

func (a *Analyzer) incr(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.stats["TotalMessages"]++
	a.stats[key]++
}

func (a *Analyzer) add(key string, val int64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.stats[key] += val
}

type AnalyzerOpts func(*Analyzer)

func WithPrintInterval(t time.Duration) AnalyzerOpts {
	return func(a *Analyzer) {
		a.printInterval = t
	}
}
