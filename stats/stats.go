package stats
import (
	"github.com/gizak/termui"
	"github.com/cloudfoundry/cli/cf/terminal"
	"time"
	"github.com/wfernandes/firehose-stats/firehose"
	"github.com/wfernandes/firehose-stats/charts"
	"github.com/cloudfoundry/cli/plugin"
)


type Stats struct {
	client  *firehose.Client
	cfUI     terminal.UI
	cliConn plugin.CliConnection
}

func New(firehoseClient *firehose.Client, cui terminal.UI, cli plugin.CliConnection) *Stats {
	return &Stats{
		client: firehoseClient,
		cfUI: cui,
		cliConn: cli,
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

		sinkTypeChart := &charts.SinkTypeChart{}
		sinkTypeChart.Init(s.cfUI)

		uaaChart := &charts.UAAChart{}
		uaaChart.Init(s.cfUI)

		msgLossChart := &charts.MsgLossChart{}
		msgLossChart.Init(s.cfUI)

		notesChart := &charts.NotesChart{}
		notesChart.Init()


		s.client.Sift(
			[]charts.Chart{
				sinkTypeChart,
				uaaChart,
				msgLossChart,
			},
		)

		termui.Body.AddRows(
			termui.NewRow(
				termui.NewCol(6, 0, sinkTypeChart),
				termui.NewCol(6, 0, uaaChart),
			),
			termui.NewRow(
				termui.NewCol(6, 0, msgLossChart),
				termui.NewCol(6, 0, notesChart),
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



