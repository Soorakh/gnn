package input

import (
	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
)

func searchEsc(ch string, s *state.State) {
	s.Search.IsActive = false
	output.UpdateScreen(s)
}

func searchBackspace(ch string, s *state.State) {
	if len(s.Search.Keyword) == 0 {
		return
	}
	s.Search.Keyword = s.Search.Keyword[0 : len(s.Search.Keyword)-1]
	files.UpdateDir(s.Dir, s, true)
	output.UpdateScreen(s)
}

func searchAppend(ch string, s *state.State) {
	s.Search.Keyword = s.Search.Keyword + ch
	files.UpdateDir(s.Dir, s, true)
	output.UpdateScreen(s)
}
