package main
import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"os"
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
				Alias: "firestats",
				HelpText: "Displays real time statistics from the Firehose. Must be logged in as an admin user.",
				UsageDetails: plugin.Usage{
					Usage: "cf firehose-stats",
					Options: map[string]string {
						"debug": "-d, enable debugging",
					},
				},
			},
		},
	}
}

func (s *FirehoseStatsCmd) Run( cliConnection plugin.CliConnection, args []string) {

	if args[0] != "firehose-stats" {
		return
	}

	s.cfUI = terminal.NewUI(os.Stdin, terminal.NewTeePrinter())

	s.cfUI.Say("This is the firehose stats command")
}