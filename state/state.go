package state

import "os"

type prev struct {
	Dir   string
	File  os.FileInfo
	Index int
}

type selected struct {
	File  os.FileInfo
	Index int
}

type Input struct {
	IsActive bool
	Keyword  string
}

type State struct {
	Dir        string
	Files      []os.FileInfo
	Selected   selected
	ShowHidden bool
	Prev       prev
	Search     Input
	Message    string
	IsPromting bool
	Rename     Input
	Mkdir      Input
}

func CreateState() *State {
	return &State{}
}
