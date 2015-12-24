package charts

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gizak/termui"
)

type UAAChart struct {
	graph            *termui.BarChart
	validOrigins     []string
	validMetricNames []string
	data             []int
	dataByIp         []map[string]int
	cfUI             terminal.UI
}

func (u *UAAChart) Init(ui terminal.UI) {
	u.data = make([]int, 4)

	u.dataByIp = make([]map[string]int, 4)
	for i := 0; i < 4; i++ {
		u.dataByIp[i] = make(map[string]int)
	}

	u.cfUI = ui
	u.graph = termui.NewBarChart()
	u.graph.BorderLabel = "UAA Metrics"
	u.graph.Data = u.data
	u.graph.Width = 80
	u.graph.Height = 20
	u.graph.DataLabels = []string{"AuthSuccess", "AuthFailure", "PrnAuthFailure", "PwdFailure"}
	u.graph.TextColor = termui.ColorYellow
	u.graph.BarColor = termui.ColorCyan
	u.graph.NumColor = termui.ColorBlack
	u.graph.BarWidth = 17
	//	u.graph.BarGap = 5

	u.validOrigins = []string{"uaa"}
	u.validMetricNames = []string{
		"audit_service.user_authentication_count",
		"audit_service.user_authentication_failure_count",
		"audit_service.principal_authentication_failure_count",
		"audit_service.user_password_failures",
	}
}

func (s *UAAChart) ForChart(event *events.Envelope) bool {
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

func (m *UAAChart) Buffer() termui.Buffer {
	return m.graph.Buffer()
}

func (m *UAAChart) GetHeight() int {
	return m.graph.GetHeight()
}

func (m *UAAChart) SetWidth(w int) {
	m.graph.SetWidth(w)
}

func (m *UAAChart) SetX(x int) {
	m.graph.SetX(x)
}

func (m *UAAChart) SetY(y int) {
	m.graph.SetY(y)
}

func (s *UAAChart) ProcessEvent(event *events.Envelope) {
	switch event.GetValueMetric().GetName() {
	case "audit_service.user_authentication_count":
		s.data[0] = updateAndReturnValue(s.dataByIp[0], event)
	case "audit_service.user_authentication_failure_count":
		s.data[1] = updateAndReturnValue(s.dataByIp[1], event)
	case "audit_service.principal_authentication_failure_count":
		s.data[2] = updateAndReturnValue(s.dataByIp[2], event)
	case "audit_service.user_password_failures":
		s.data[3] = updateAndReturnValue(s.dataByIp[3], event)
	}
}
