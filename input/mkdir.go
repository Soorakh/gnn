package input

import (
	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
)

func mkdirEsc(ch string, s *state.State) {
	s.Mkdir.IsActive = false
	s.Mkdir.Keyword = ""
	output.UpdateScreen(s)
}

func mkdirEnter(ch string, s *state.State) {
	err := files.CreateDirectory(s.Mkdir.Keyword)

	if err != nil {
		s.Message = err.Error()
	}

	s.Mkdir.IsActive = false
	s.Mkdir.Keyword = ""

	// todo select created dir
	files.UpdateDir(s.Dir, s, false)
	output.UpdateScreen(s)
}

func mkdirBackspace(ch string, s *state.State) {
	if len(s.Mkdir.Keyword) == 0 {
		return
	}
	s.Mkdir.Keyword = s.Mkdir.Keyword[0 : len(s.Mkdir.Keyword)-1]
	output.UpdateScreen(s)
}

func mkdirAppend(ch string, s *state.State) {
	s.Mkdir.Keyword = s.Mkdir.Keyword + ch
	output.UpdateScreen(s)
}
