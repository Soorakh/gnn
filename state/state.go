package state

import "os"

type prev struct {
	Dir   string
	File  os.FileInfo
	Index int
}

type State struct {
	Dir               string
	Files             []os.FileInfo
	SelectedFile      os.FileInfo
	SelectedFileIndex int
	ShowHidden        bool
	Prev              prev
}

func CreateState() *State {
	return &State{}
}
