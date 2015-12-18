package firehose
import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/wfernandes/firehose-stats/charts"
)

type Sifter interface {
	Sift(in <-chan *events.Envelope, charts []charts.Chart)
}

type ChartFilter struct {

}


func (of *ChartFilter) Sift(in <-chan *events.Envelope, charts []charts.Chart) {
	go func() {
		for inEvent := range in {
			for _, chart := range charts {
				if chart.ForChart(inEvent) {
					chart.ProcessEvent(inEvent)
				}
			}
		}
	}()
}
