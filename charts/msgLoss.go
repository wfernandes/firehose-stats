package charts
import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gizak/termui"
)

type Chart interface {
	ProcessEvent(e *events.Envelope)
	ForChart(e *events.Envelope) bool
	Buffer() termui.Buffer
}


type MsgLossChart struct {
	graph *termui.Gauge
	validOrigins []string
	validMetricNames []string
}

func (m *MsgLossChart) Init() {
	m.graph = termui.NewGauge()
	m.graph.Width = 50
	m.graph.Height = 3
	m.graph.PercentColor = termui.ColorBlue
	m.graph.Y = 0
	m.graph.X = 0
	m.graph.BorderLabel = "Slim Gauge"
	m.graph.BarColor = termui.ColorYellow
	m.graph.BorderFg = termui.ColorWhite


	m.validOrigins = []string{"MetronAgent", "DopplerServer"}
	m.validMetricNames = []string{"tls.sentMessageCount", "tlsListener.receivedMessageCount"}

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

func (m* MsgLossChart) Buffer() termui.Buffer {
	return m.graph.Buffer()
}


func (m * MsgLossChart) ProcessEvent() {
	// do nothing
}


func contains(name string, list []string) bool {

	for _, lv := range list {
		if lv == name {
			return true
		}
	}
	return false
}