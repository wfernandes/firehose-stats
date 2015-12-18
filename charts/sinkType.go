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
	cfUI			terminal.UI
}

func (s *SinkTypeChart) Init(ui terminal.UI) {
	s.data = make([]int, 5)
	s.cfUI = ui
	s.graph = termui.NewBarChart()
	s.graph.BorderLabel = "Bar Chart"
	s.graph.Data = s.data
	s.graph.Width = 100
	s.graph.Height = 10
	s.graph.DataLabels = []string{"ContainerMetric", "SysLog", "Dump", "Websocket", "Firehose"}
	s.graph.TextColor = termui.ColorGreen
	s.graph.BarColor = termui.ColorRed
	s.graph.NumColor = termui.ColorYellow
	s.graph.BarWidth = 15


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


func (s * SinkTypeChart) ProcessEvent(event *events.Envelope) {
	switch event.GetValueMetric().GetName() {
	case "messageRouter.numberOfContainerMetricSinks":
		s.data[0] = int(event.GetValueMetric().GetValue())
	case "messageRouter.numberOfSyslogSinks":
		s.data[1] = int(event.GetValueMetric().GetValue())
	case "messageRouter.numberOfDumpSinks":
		s.data[2] = int(event.GetValueMetric().GetValue())
	case "messageRouter.numberOfWebsocketSinks":
		s.data[3] = int(event.GetValueMetric().GetValue())
	case "messageRouter.numberOfFirehoseSinks":
		s.data[4] = int(event.GetValueMetric().GetValue())
	}
}
