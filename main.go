package main

import (
	"github.com/Soorakh/gnn/events"
	"github.com/Soorakh/gnn/state"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	termbox.SetOutputMode(termbox.Output256)
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	s := state.CreateState()
	events.Init(s)
	events.Bind(s)
}
