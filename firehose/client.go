package firehose

import (
	"crypto/tls"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
)

//
//type Sifter interface {
//	Sift(charts []charts.Chart)
//}

type Consumer interface {
	Consume(e *events.Envelope)
}

type Client struct {
	dopplerEndpoint string
	authToken       string
	ui              terminal.UI
	outputChan      chan *events.Envelope
	consumer        Consumer
}

func NewClient(authToken, doppplerEndpoint string, consumer Consumer, ui terminal.UI) *Client {
	return &Client{
		dopplerEndpoint: doppplerEndpoint,
		authToken:       authToken,
		ui:              ui,
		consumer:        consumer,
	}

}

func (c *Client) Start() {

	firehoseChan := make(chan *events.Envelope)
	dopplerConnection := noaa.NewConsumer(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	subscriptionID := "firehose-stats"
	go func() {
		err := dopplerConnection.FirehoseWithoutReconnect(subscriptionID, c.authToken, firehoseChan)
		if err != nil {
			c.ui.Warn(err.Error())
			close(firehoseChan)
			return
		}
	}()

	defer dopplerConnection.Close()
	c.ui.Say("Starting the nozzle")

	for envelope := range firehoseChan {
		c.consumer.Consume(envelope)
	}

}

//
//func (c *Client) Sift(charts []charts.Chart) {
//	go func() {
//		for inEvent := range c.outputChan {
//			for _, chart := range charts {
//				if chart.ForChart(inEvent) {
//					chart.ProcessEvent(inEvent)
//				}
//			}
//		}
//	}()
//}
