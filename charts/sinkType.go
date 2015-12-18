package charts
import (
	"github.com/gizak/termui"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/cloudfoundry/cli/cf/terminal"
)


type SinkTypeChart struct {
	graph            *termui.BarChart
	validOrigins     []string
	validMetricNames []string
	data             []int
	dataByIp		 []map[string]int
	cfUI			terminal.UI
}

func (s *SinkTypeChart) Init(ui terminal.UI) {
	s.data = make([]int, 5)

	s.dataByIp = make([]map[string]int, 5)
	for i :=0 ; i < 5; i++ {
		s.dataByIp[i] = make(map[string]int)
	}

	s.cfUI = ui
	s.graph = termui.NewBarChart()
	s.graph.BorderLabel = "Number of Sinks"
	s.graph.Data = s.data
	s.graph.Width = 80
	s.graph.Height = 20
	s.graph.DataLabels = []string{"ContainerMetric", "SysLog", "Dump", "Websocket", "Firehose"}
	s.graph.TextColor = termui.ColorGreen
	s.graph.BarColor = termui.ColorRed
	s.graph.NumColor = termui.ColorYellow
	s.graph.BarWidth = 10


	s.validOrigins = []string{"DopplerServer"}
	s.validMetricNames = []string{
		"messageRouter.numberOfContainerMetricSinks",
		"messageRouter.numberOfSyslogSinks",
		"messageRouter.numberOfDumpSinks",
		"messageRouter.numberOfWebsocketSinks",
		"messageRouter.numberOfFirehoseSinks",
	}
}

func (s *SinkTypeChart) ForChart(event *events.Envelope) bool {
	if event.GetEventType() != events.Envelope_ValueMetric {
		return false
	}
	if !contains(event.GetOrigin(), s.validOrigins) {
		return false
	}
	if !contains(event.GetValueMetric().GetName(), s.validMetricNames) {
		return false
	}

//	s.cfUI.Say("%f ", event.GetValueMetric().GetValue())
	return true
}

func (m* SinkTypeChart) Buffer() termui.Buffer {
	return m.graph.Buffer()
}

func (m* SinkTypeChart) GetHeight() int {
	return m.graph.GetHeight()
}

func (m* SinkTypeChart) SetWidth(w int) {
	m.graph.SetWidth(w)
}

func (m* SinkTypeChart) SetX(x int) {
	m.graph.SetX(x)
}

func (m* SinkTypeChart) SetY(y int) {
	m.graph.SetY(y)
}


func (s * SinkTypeChart) ProcessEvent(event *events.Envelope) {
	switch event.GetValueMetric().GetName() {
	case "messageRouter.numberOfContainerMetricSinks":
		s.data[0] = updateAndReturnValue(s.dataByIp[0], event)
	case "messageRouter.numberOfSyslogSinks":
		s.data[1] = updateAndReturnValue(s.dataByIp[1], event)
	case "messageRouter.numberOfDumpSinks":
		s.data[2] = updateAndReturnValue(s.dataByIp[2], event)
	case "messageRouter.numberOfWebsocketSinks":
		s.data[3] = updateAndReturnValue(s.dataByIp[3], event)
	case "messageRouter.numberOfFirehoseSinks":
		s.data[4] = updateAndReturnValue(s.dataByIp[4], event)
	}
}

func updateAndReturnValue(values map[string]int, event *events.Envelope) int {
	values[event.GetIp()] = int(event.GetValueMetric().GetValue())

	sum := 0
	for _, v := range values {
		sum += v
	}

	return sum
}
