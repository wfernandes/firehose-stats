package firehose

import (
	"crypto/tls"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
)

type Client struct {
	dopplerEndpoint string
	authToken       string
	ui              terminal.UI
	outputChan		chan *events.Envelope
}

func NewClient(authToken, doppplerEndpoint string, ui terminal.UI, outChan chan *events.Envelope) *Client {
	return &Client{
		dopplerEndpoint: doppplerEndpoint,
		authToken:       authToken,
		ui:              ui,
		outputChan: outChan,
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
