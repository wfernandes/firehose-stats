package main

import (
	"fmt"
	"sort"
	"strings"

	tm "github.com/buger/goterm"
)

type ResultsPrinter struct{}

func (p *ResultsPrinter) Print(s Stats) {
	// fmt.Printf("%#v\n", s)
	tm.Clear()
	tm.MoveCursor(1, 1)

	var headers []string
	for k, _ := range s {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var values []string
	for _, h := range headers {
		values = append(values, fmt.Sprintf("%d", s[h]))
	}

	output := tm.NewTable(0, 20, 5, ' ', 0)
	fmt.Fprintf(output, fmt.Sprintf("%s\n", strings.Join(headers, "\t")))
	fmt.Fprintf(output, fmt.Sprintf("%s\n", strings.Join(values, "\t")))
	tm.Println(output)
	tm.Flush()
}
