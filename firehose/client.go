package firehose

import (
	"crypto/tls"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/wfernandes/firehose-stats/charts"
)

type Sifter interface {
	Sift(in <-chan *events.Envelope, charts []charts.Chart)
}

type Client struct {
	dopplerEndpoint string
	authToken       string
	ui              terminal.UI
	outputChan		chan *events.Envelope
}

func NewClient(authToken, doppplerEndpoint string, ui terminal.UI) *Client {
	firehoseChan := make(chan *events.Envelope)
	return &Client{
		dopplerEndpoint: doppplerEndpoint,
		authToken:       authToken,
		ui:              ui,
		outputChan: firehoseChan,
	}

}

func (c *Client) Start() {

	dopplerConnection := noaa.NewConsumer(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	subscriptionID := "firehose-stats"
	go func() {
		err := dopplerConnection.FirehoseWithoutReconnect(subscriptionID, c.authToken, c.outputChan)
		if err != nil {
			c.ui.Warn(err.Error())
			close(c.outputChan)
			return
		}
	}()

	defer dopplerConnection.Close()

	c.ui.Say("Starting the nozzle")
	c.ui.Say("Hit Ctrl+c to exit")

}


func (c *Client) Sift(charts []charts.Chart) {
	go func() {
		for inEvent := range c.outputChan {
			for _, chart := range charts {
				if chart.ForChart(inEvent) {
					chart.ProcessEvent(inEvent)
				}
			}
		}
	}()
}
