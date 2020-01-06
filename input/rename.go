package input

import (
	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
)

func renameEsc(ch string, s *state.State) {
	s.Rename.IsActive = false
	s.Rename.Keyword = ""
	output.UpdateScreen(s)
}

func renameEnter(ch string, s *state.State) {
	err := files.MoveFile(s.Dir, s.Selected.File, s.Rename.Keyword)

	if err != nil {
		s.Message = err.Error()
	}

	s.Rename.IsActive = false
	s.Rename.Keyword = ""

	files.UpdateDir(s.Dir, s, false)
	output.UpdateScreen(s)
}

func renameBackspace(ch string, s *state.State) {
	if len(s.Rename.Keyword) == 0 {
		return
	}
	s.Rename.Keyword = s.Rename.Keyword[0 : len(s.Rename.Keyword)-1]
	output.UpdateScreen(s)
}

func renameAppend(ch string, s *state.State) {
	s.Rename.Keyword = s.Rename.Keyword + ch
	output.UpdateScreen(s)
}
