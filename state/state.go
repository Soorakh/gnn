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

type State struct {
	Dir        string
	Files      []os.FileInfo
	Selected   selected
	ShowHidden bool
	Prev       prev
}

func CreateState() *State {
	return &State{}
}
