package main

import (
	"fmt"
	"os"
)

type Logger struct {
	Verbose bool
}

func (l Logger) Debug(a ...interface{}) {
	if l.Verbose {
		fmt.Println(a)
	}
}

func (l Logger) Info(a ...interface{}) {
	fmt.Println(a)
}

func (l Logger) Fatal(a ...interface{}) {
	fmt.Println(a)
	os.Exit(1)
}
