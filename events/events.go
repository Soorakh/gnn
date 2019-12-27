package events

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Soorakh/gnn/cursor"
	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
	"github.com/nsf/termbox-go"
)

func Init(s *state.State) {
	dir, _ := os.Getwd()
	updateDir(dir, s, true)
	updateScreen(s)
}

func Bind(s *state.State) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Ch == 'q':
				return
			case ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown:
				moveCursorDown(s)
			case ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp:
				moveCursorUp(s)
			case ev.Ch == 'l' || ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyEnter:
				openFile(s.SelectedFile, s)
			case ev.Ch == 'h' || ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyBackspace2:
				changeDirUp(s)
			case ev.Ch == '.':
				toggleHidden(s)
			case ev.Ch == 'e':
				editFile(s)
			}
		case termbox.EventResize:
			updateScreen(s)
		}
	}
}

func moveCursorDown(s *state.State) {
	s.SelectedFileIndex = cursor.MoveDown(s.SelectedFileIndex, len(s.Files))
	s.SelectedFile = s.Files[s.SelectedFileIndex]
	updateScreen(s)
}

func updateScreen(s *state.State) {
	output.PrintFiles(s.Files, s.SelectedFile, s.Dir, s.SelectedFileIndex)
}

func moveCursorUp(s *state.State) {
	s.SelectedFileIndex = cursor.MoveUp(s.SelectedFileIndex, len(s.Files))
	s.SelectedFile = s.Files[s.SelectedFileIndex]
	updateScreen(s)
}

func openFile(file os.FileInfo, s *state.State) {
	if file == nil {
		return
	}
	if file.IsDir() {
		delimeter := "/"
		if s.Dir == "/" {
			delimeter = ""
		}
		s.Prev.Dir = s.Dir
		s.Prev.File = file
		updateDir(s.Dir+delimeter+file.Name(), s, true)
		updateScreen(s)
	} else {
		cmd := exec.Command("xdg-open", s.Dir+"/"+file.Name())
		cmd.Start()
	}
}

func updateDir(d string, s *state.State, resetSelected bool) {
	ratio := float64(s.SelectedFileIndex) / float64(len(s.Files))

	s.Dir = d
	s.Files = files.GetFiles(d, s.ShowHidden)
	if !resetSelected {
		if s.SelectedFileIndex >= len(s.Files) {
			s.SelectedFileIndex = int(ratio * float64(len(s.Files)))
		}
		for i, v := range s.Files {
			if v.Name() == s.SelectedFile.Name() {
				s.SelectedFileIndex = i
				break
			}
		}
	} else {
		s.SelectedFileIndex = 0
	}
	if len(s.Files) > 0 {
		s.SelectedFile = s.Files[s.SelectedFileIndex]
	} else {
		s.SelectedFile = nil
	}
}

func changeDirUp(s *state.State) {
	p := strings.Split(s.Dir, "/")
	plen := len(p)
	if plen < 2 {
		return
	}
	p = p[:plen-1]
	newDir := "/"
	if plen > 2 {
		newDir = strings.Join(p, "/")
	}

	resetSelected := true

	if newDir == s.Prev.Dir {
		s.SelectedFileIndex = s.Prev.Index
		s.SelectedFile = s.Prev.File
		resetSelected = false
	}

	updateDir(newDir, s, resetSelected)

	updateScreen(s)
}

func toggleHidden(s *state.State) {
	s.ShowHidden = !s.ShowHidden
	updateDir(s.Dir, s, false)
	updateScreen(s)
}

func editFile(s *state.State) {
	editor := os.Getenv("EDITOR")
	cmd := exec.Command(editor, s.SelectedFile.Name())
	cmd.Dir = s.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	output.FixScreen()
	updateScreen(s)
}
