package main

import "fmt"

type ResultsPrinter struct{}

func (p *ResultsPrinter) Print(s Stats) {
	fmt.Printf("%#v\n", s)
}
