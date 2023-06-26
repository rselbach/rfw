package main

import "fmt"

type logger bool

func (l logger) log(a ...any) {
	if l {
		fmt.Println(a...)
	}
}

func (l logger) logf(tmpl string, a ...any) {
	if l {
		fmt.Printf(tmpl, a...)
		fmt.Println()
	}
}
