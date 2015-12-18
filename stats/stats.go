package stats
import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gizak/termui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"time"
	"github.com/wfernandes/firehose-stats/firehose"
	"github.com/wfernandes/firehose-stats/charts"
)


type Stats struct {
	dataChan <-chan *events.Envelope
	cfUI     terminal.UI
}

func New(output chan *events.Envelope, cui terminal.UI) *Stats {
	return &Stats{
		dataChan: output,
		cfUI: cui,
	}
}

func (s *Stats) Start() {
	s.cfUI.Say("Starting Stats...")
	err := termui.Init()
	if err != nil {
		s.cfUI.Warn(err.Error())
		return
	}
	defer termui.Close()

	go func() {

		firehoseMF := &firehose.ChartFilter{}

		//
		//		msgLossChart := &charts.MsgLossChart{}
		//		msgLossChart.Init()


		sinkTypeChart := &charts.SinkTypeChart{}
		sinkTypeChart.Init(s.cfUI)

		uaaChart := &charts.UAAChart{}
		uaaChart.Init(s.cfUI)

		firehoseMF.Sift(
			s.dataChan,
			[]charts.Chart{
				sinkTypeChart,
				uaaChart,
			},
		)


		termui.Body.AddRows(
			termui.NewRow(
				termui.NewCol(6, 0, sinkTypeChart),
			),
			termui.NewRow(
				termui.NewCol(6, 0, uaaChart),
			),
		)



		for {
			termui.Body.Align()
			termui.Render(termui.Body)
			time.Sleep(1 * time.Second)
		}
	}()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()

}



