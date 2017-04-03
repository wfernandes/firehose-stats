package charts

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gizak/termui"
)

type Chart interface {
	termui.GridBufferer
	ForChart(e *events.Envelope) bool
	ProcessEvent(e *events.Envelope)
}

type MsgLossChart struct {
	graph *termui.Gauge
	cfUI  terminal.UI

	validOrigins     []string
	validMetricNames []string
	totalSent        int64
	totalReceived    int64

	sentByIP     map[string]int64
	receivedByIP map[string]int64
}

func (m *MsgLossChart) Init(cfUI terminal.UI) {
	m.sentByIP = make(map[string]int64)
	m.receivedByIP = make(map[string]int64)
	m.cfUI = cfUI
	m.graph = termui.NewGauge()
	m.graph.Width = 50
	m.graph.Height = 8
	m.graph.PercentColor = termui.ColorBlue
	m.graph.Y = 0
	m.graph.X = 0
	m.graph.BorderLabel = "(%)Msg Loss Between Metron and Doppler"
	m.graph.BarColor = termui.ColorYellow
	m.graph.BorderFg = termui.ColorWhite

	m.validOrigins = []string{"MetronAgent", "DopplerServer"}
	m.validMetricNames = []string{"DopplerForwarder.sentMessages", "tlsListener.receivedMessageCount", "dropsondeListener.receivedMessageCount"}

}

func (m *MsgLossChart) ForChart(event *events.Envelope) bool {
	if event.GetEventType() != events.Envelope_CounterEvent {
		return false
	}
	if !contains(event.GetOrigin(), m.validOrigins) {
		return false
	}
	if !contains(event.GetCounterEvent().GetName(), m.validMetricNames) {
		return false
	}
	return true
}

func (m *MsgLossChart) Buffer() termui.Buffer {
	return m.graph.Buffer()
}

func (m *MsgLossChart) GetHeight() int {
	return m.graph.GetHeight()
}

func (m *MsgLossChart) SetWidth(w int) {
	m.graph.SetWidth(w)
}

func (m *MsgLossChart) SetX(x int) {
	m.graph.SetX(x)
}

func (m *MsgLossChart) SetY(y int) {
	m.graph.SetY(y)
}

func (m *MsgLossChart) ProcessEvent(evt *events.Envelope) {
	switch evt.GetCounterEvent().GetName() {
	case "DopplerForwarder.sentMessages":
		m.totalSent = update(m.sentByIP, evt)
	case "tlsListener.receivedMessageCount", "dropsondeListener.receivedMessageCount":
		m.totalReceived = update(m.receivedByIP, evt)
	}
	percent := 100 * ((m.totalSent - m.totalReceived) / m.totalSent)

	m.graph.Percent = int(percent)
}

func contains(name string, list []string) bool {

	for _, lv := range list {
		if lv == name {
			return true
		}
	}
	return false
}

func update(values map[string]int64, event *events.Envelope) int64 {

	values[event.GetIp()] = int64(event.GetCounterEvent().GetTotal())

	var sum int64
	for _, v := range values {
		sum = sum + v
	}

	return sum
}
