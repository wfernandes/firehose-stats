package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/noaa/consumer"
)

type FirehoseStatsCmd struct {
	ui terminal.UI
}

func main() {
	plugin.Start(new(FirehoseStatsCmd))
}

func (s *FirehoseStatsCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "FirehoseStats",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 2,
		},
		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 3,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "firehose-stats",
				Alias:    "fs",
				HelpText: "Displays real time statistics from the Firehose. Must be logged in as an admin user.",
				UsageDetails: plugin.Usage{
					Usage: "cf firehose-stats",
					Options: map[string]string{
						"debug": "-d, enable debugging",
					},
				},
			},
		},
	}
}

func (s *FirehoseStatsCmd) Run(cliConnection plugin.CliConnection, args []string) {

	if args[0] != "firehose-stats" {
		return
	}

	s.ui = terminal.NewUI(os.Stdin, terminal.NewTeePrinter())

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		s.ui.Failed(err.Error())
	}

	authToken, err := cliConnection.AccessToken()
	if err != nil {
		s.ui.Failed(err.Error())
	}

	consumer := consumer.New(dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	defer consumer.Close()

	// Create a unique subscription id to avoid collisions
	subscriptionID := fmt.Sprintf("firehose-stats-%d", time.Now().UnixNano())
	msgs, errs := consumer.Firehose(subscriptionID, authToken)

	go func() {
		for e := range errs {
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
	}()

	s.ui.Say("Starting the nozzle")
	s.ui.Say("Hit Ctrl+c to exit")

	p := &ResultsPrinter{}

	analyzer := NewAnalyzer(msgs, p)
	analyzer.Start()
}
