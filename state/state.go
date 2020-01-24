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
	Offset   int
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
	C          chan State
}

func CreateState() *State {
	s := &State{}
	s.C = make(chan State)
	return s
}

func (s *State) Apply() {
	s.C <- *s
}
