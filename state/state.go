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

type Search struct {
	IsActive bool
	Keyword  string
}

type State struct {
	Dir        string
	Files      []os.FileInfo
	Selected   selected
	ShowHidden bool
	Prev       prev
	Search     Search
	Message    string
	IsPromting bool
}

func CreateState() *State {
	return &State{}
}
