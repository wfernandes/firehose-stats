package main

import (
	"sync/atomic"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

type Printer interface {
	Print(Stats)
}

type Stats map[string]int64

var GlobalStats atomic.Value

type Analyzer struct {
	messages      <-chan *events.Envelope
	printer       Printer
	printInterval time.Duration
	stats         Stats
}

func NewAnalyzer(msg <-chan *events.Envelope, p Printer) *Analyzer {
	return &Analyzer{
		messages:      msg,
		printer:       p,
		printInterval: 5 * time.Second,
		stats:         Stats{},
	}
}

func (a *Analyzer) Start() {

	ticker := time.NewTicker(a.printInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			data := GlobalStats.Load().(Stats)
			a.printer.Print(data)
		}
	}()

	for e := range a.messages {
		a.stats["TotalMessages"]++
		a.stats["AvgEnvelopeSize"] = a.stats["AvgEnvelopeSize"] + int64(e.Size())/a.stats["TotalMessages"]
		switch e.GetEventType() {
		case events.Envelope_CounterEvent:
			a.stats["CounterEvents"]++
		case events.Envelope_LogMessage:
			a.stats["LogMessages"]++
		case events.Envelope_ContainerMetric:
			a.stats["ContainerMetrics"]++
		case events.Envelope_ValueMetric:
			a.stats["ValueMetrics"]++
		}
		GlobalStats.Store(a.stats)
	}
}
