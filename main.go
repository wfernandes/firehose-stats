package main
import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"github.com/wfernandes/firehose-stats/firehose"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/wfernandes/firehose-stats/stats"
)

type FirehoseStatsCmd struct {
	cfUI terminal.UI
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
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 3,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name: "firehose-stats",
				Alias: "fs",
				HelpText: "Displays real time statistics from the Firehose. Must be logged in as an admin user.",
				UsageDetails: plugin.Usage{
					Usage: "cf firehose-stats",
					Options: map[string]string{
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

	s.cfUI = terminal.NewUI(os.Stdin, terminal.NewTeePrinter())

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		s.cfUI.Failed(err.Error())
	}

	authToken, err := cliConnection.AccessToken()
	if err != nil {
		s.cfUI.Failed(err.Error())
	}
	client := firehose.NewClient(authToken, dopplerEndpoint, s.cfUI, )
	client.Start()

	statsUI := stats.New(client, s.cfUI, cliConnection)
	statsUI.Start()

}