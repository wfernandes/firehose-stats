package charts
import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gizak/termui"
)

type NotesChart struct {
	graph *termui.Par
	validOrigins []string
	validMetricNames []string
	totalSent int
	totalReceived int

	sentByIP map[string]int
	receivedByIP map[string]int
}

func (m *NotesChart) Init() {
	m.sentByIP = make(map[string]int)
	m.receivedByIP = make(map[string]int)

	m.graph = termui.NewPar("Firehose Statistics Plugin\n [Press q](fg-red) to exit.")

//	m.graph.Width = 50
	m.graph.Height = 7
	m.graph.BorderLabel = "Notes"
	m.graph.BorderFg = termui.ColorYellow

}

func (m *NotesChart) ForChart(event *events.Envelope) bool {
	return false
}

func (m* NotesChart) Buffer() termui.Buffer {
	return m.graph.Buffer()
}

func (m* NotesChart) GetHeight() int {
	return m.graph.GetHeight()
}

func (m* NotesChart) SetWidth(w int) {
	m.graph.SetWidth(w)
}

func (m* NotesChart) SetX(x int) {
	m.graph.SetX(x)
}

func (m* NotesChart) SetY(y int) {
	m.graph.SetY(y)
}

func (m * NotesChart) ProcessEvent(evt *events.Envelope) {
	// do nothing
}
