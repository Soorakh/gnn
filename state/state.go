package state

import "os"

type State struct {
	Dir               string
	Files             []os.FileInfo
	SelectedFile      os.FileInfo
	SelectedFileIndex int
	ShowHidden        bool
}

func CreateState() *State {
	return &State{}
}
