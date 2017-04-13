package main

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/olekukonko/tablewriter"
)

type ResultsPrinter struct{}

func (p *ResultsPrinter) Print(s Stats) {

	tm.Clear()
	tm.MoveCursor(1, 1)
	table := tablewriter.NewWriter(tm.Screen)
	table.SetHeader([]string{"Name", "Value"})
	for k, v := range s {
		table.Append([]string{k, fmt.Sprintf("%d", v)})
	}
	table.Render()
	tm.Flush()
}
