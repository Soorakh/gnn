package input

import (
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

func HandleMkdir(ch string, key termbox.Key, s *state.State) {
	inputHandler(ch, key, s, mkdirEsc, mkdirBackspace, mkdirAppend, mkdirEnter)
}
