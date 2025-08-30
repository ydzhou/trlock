package main

import (
	trlock "github.com/ydzhou/trlock/internal"
)

func main() {
	t := &trlock.Trlock{}

	t.Setup()
	t.Run()
}
