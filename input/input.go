package input

import (
	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
	"github.com/nsf/termbox-go"
)

type inputCallback func(s *state.State)

func inputHandler(ch string, key termbox.Key, s *state.State, input *state.Input, rescan bool, onEnter inputCallback) {
	if key == termbox.KeyEsc {
		termbox.HideCursor()
		input.IsActive = false
		output.UpdateScreen(s)
		return
	}

	if key == termbox.KeyBackspace2 {
		if input.Offset == 0 || len(input.Keyword) == 0 {
			return
		}
		input.Keyword = input.Keyword[0:input.Offset-1] + input.Keyword[input.Offset:]

		input.Offset = input.Offset - 1
		if rescan {
			files.UpdateDir(s.Dir, s, true)
		}
		output.UpdateScreen(s)
		return
	}

	if key == termbox.KeyEnter {
		termbox.HideCursor()
		onEnter(s)

		input.IsActive = false
		input.Keyword = ""

		files.UpdateDir(s.Dir, s, false)
		output.UpdateScreen(s)
		return
	}

	if key == termbox.KeyArrowLeft {
		if input.Offset == 0 {
			return
		}
		input.Offset = input.Offset - 1
		output.UpdateScreen(s)
	}

	if key == termbox.KeyArrowRight {
		if input.Offset == len(input.Keyword) {
			return
		}
		input.Offset = input.Offset + 1
		output.UpdateScreen(s)
	}

	if key == termbox.KeyDelete {
		if input.Offset == len(input.Keyword) {
			return
		}
		input.Keyword = input.Keyword[:input.Offset] + input.Keyword[input.Offset+1:]
		if rescan {
			files.UpdateDir(s.Dir, s, false)
		}
		output.UpdateScreen(s)
	}

	// Non-iput buttons are ignored
	if key <= termbox.KeyF1 && key >= termbox.KeyArrowRight {
		return
	}

	input.Keyword = input.Keyword[:input.Offset] + ch + input.Keyword[input.Offset:]
	if rescan {
		files.UpdateDir(s.Dir, s, true)
	}
	input.Offset = input.Offset + 1
	output.UpdateScreen(s)
}

func HandleSearch(ch string, key termbox.Key, s *state.State) {
	inputHandler(ch, key, s, &s.Search, true, func(s *state.State) {})
}

func HandleRename(ch string, key termbox.Key, s *state.State) {
	inputHandler(ch, key, s, &s.Rename, false, renameEnter)
}

func HandleMkdir(ch string, key termbox.Key, s *state.State) {
	inputHandler(ch, key, s, &s.Mkdir, false, mkdirEnter)
}

func renameEnter(s *state.State) {
	err := files.MoveFile(s.Dir, s.Selected.File, s.Rename.Keyword)

	if err != nil {
		s.Message = err.Error()
	}
}

func mkdirEnter(s *state.State) {
	// todo select created dir
	err := files.CreateDirectory(s.Mkdir.Keyword)

	if err != nil {
		s.Message = err.Error()
	}
}
