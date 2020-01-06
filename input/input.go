package input

import (
	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
	"github.com/nsf/termbox-go"
)

type inputCallback func(ch string, s *state.State)

func inputHandler(ch string, key termbox.Key, s *state.State, esc inputCallback, backspace inputCallback, append inputCallback, enter inputCallback) {
	if key == termbox.KeyEsc {
		esc("", s)
		return
	}

	if key == termbox.KeyBackspace2 {
		backspace("", s)
		return
	}

	if key == termbox.KeyEnter {
		enter("", s)
		return
	}

	// Non-iput buttons are ignored
	if key <= termbox.KeyF1 && key >= termbox.KeyArrowRight {
		return
	}

	append(ch, s)
}

func HandleSearch(ch string, key termbox.Key, s *state.State) {
	inputHandler(ch, key, s, searchEsc, searchBackspace, searchAppend, searchEsc)
}

func HandleRename(ch string, key termbox.Key, s *state.State) {
	inputHandler(ch, key, s, renameEsc, renameBackspace, renameAppend, renameEnter)
}

func searchEsc(ch string, s *state.State) {
	s.Search.IsActive = false
	output.UpdateScreen(s)
}

func renameEsc(ch string, s *state.State) {
	s.Rename.IsActive = false
	s.Rename.Keyword = ""
	output.UpdateScreen(s)
}

func renameEnter(ch string, s *state.State) {
	err := files.RenameFile(s.Dir, s.Selected.File, s.Rename.Keyword)

	if err != nil {
		s.Message = err.Error()
	}

	s.Rename.IsActive = false
	s.Rename.Keyword = ""

	files.UpdateDir(s.Dir, s, false)
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

func renameBackspace(ch string, s *state.State) {
	if len(s.Rename.Keyword) == 0 {
		return
	}
	s.Rename.Keyword = s.Rename.Keyword[0 : len(s.Rename.Keyword)-1]
	output.UpdateScreen(s)
}

func searchAppend(ch string, s *state.State) {
	s.Search.Keyword = s.Search.Keyword + ch
	files.UpdateDir(s.Dir, s, true)
	output.UpdateScreen(s)
}

func renameAppend(ch string, s *state.State) {
	s.Rename.Keyword = s.Rename.Keyword + ch
	output.UpdateScreen(s)
}
