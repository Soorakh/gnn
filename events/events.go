package events

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/Soorakh/gnn/files"
	"github.com/Soorakh/gnn/input"
	"github.com/Soorakh/gnn/output"
	"github.com/Soorakh/gnn/state"
	"github.com/nsf/termbox-go"
)

func Init(s *state.State) {
	dir, _ := os.Getwd()
	files.UpdateDir(dir, s, true)
	output.UpdateScreen(s)
}

func Bind(s *state.State) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch {
			case s.Search.IsActive:
				input.HandleSearch(string(ev.Ch), ev.Key, s)
			case s.Rename.IsActive:
				input.HandleRename(string(ev.Ch), ev.Key, s)
			case s.Mkdir.IsActive:
				input.HandleMkdir(string(ev.Ch), ev.Key, s)
			case s.IsPromting:
				checkPromt(ev.Ch, s)
			case ev.Key == termbox.KeyEsc && s.Search.Keyword != "":
				cancelSearch(s)
			case ev.Ch == 'q':
				return
			case ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown:
				moveCursorDown(s)
			case ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp:
				moveCursorUp(s)
			case ev.Ch == 'l' || ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyEnter:
				openFile(s.Selected.File, s)
			case ev.Ch == 'h' || ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyBackspace2:
				changeDirUp(s)
			case ev.Ch == '.':
				toggleHidden(s)
			case ev.Ch == 'e':
				editFile(s)
			case ev.Ch == '/':
				searchToggleOn(s)
			case ev.Ch == 'd':
				togglePromtOn(s)
			case ev.Ch == 'r':
				files.UpdateDir(s.Dir, s, false)
				output.UpdateScreen(s)
			case ev.Ch == 'm':
				toggleRenameOn(s)
			case ev.Ch == 'n':
				toggleMkdirOn(s)
			case ev.Ch == 'y':
				yankFile(s)
			case ev.Ch == 'p':
				putFile(s)
			}
		case termbox.EventResize:
			output.UpdateScreen(s)
		}
	}
}

func yankFile(s *state.State) {
	// copy-paste should be process independent
	// temporary file is best solution here
	tmp := os.TempDir() + "/gnn"
	f, err := os.Create(tmp)
	if err != nil {
		s.Message = err.Error()
		output.UpdateScreen(s)
		return
	}

	defer f.Close()

	_, err = f.Write([]byte(filepath.Join(s.Dir, s.Selected.File.Name())))

	if err != nil {
		s.Message = err.Error()
		output.UpdateScreen(s)
		return
	}

	s.Message = "Yanked!"
	output.UpdateScreen(s)
}

func putFile(s *state.State) {
	// TODO read by lines and make array of files to copy
	f, err := ioutil.ReadFile(os.TempDir() + "/gnn")

	if err != nil {
		s.Message = err.Error()
		output.UpdateScreen(s)
		return
	}

	err = files.CopyFile(string(f), s.Dir)

	if err != nil {
		s.Message = err.Error()
		output.UpdateScreen(s)
		return
	}

	files.UpdateDir(s.Dir, s, false)
	output.UpdateScreen(s)
}

func cancelSearch(s *state.State) {
	s.Search.Keyword = ""
	files.UpdateDir(s.Dir, s, false)
	output.UpdateScreen(s)
}

func searchToggleOn(s *state.State) {
	s.Search.Keyword = ""
	s.Search.IsActive = true
	s.Search.Offset = 0
	output.UpdateScreen(s)
}

func toggleRenameOn(s *state.State) {
	if s.Selected.File == nil {
		return
	}
	s.Rename.Keyword = filepath.Join(s.Dir, s.Selected.File.Name())
	s.Rename.IsActive = true
	s.Rename.Offset = len(s.Rename.Keyword)
	output.UpdateScreen(s)
}

func toggleMkdirOn(s *state.State) {
	trailing := string(filepath.Separator)
	if s.Dir == string(filepath.Separator) {
		trailing = ""
	}
	s.Mkdir.Keyword = s.Dir + trailing
	s.Mkdir.IsActive = true
	s.Mkdir.Offset = len(s.Mkdir.Keyword)
	output.UpdateScreen(s)
}

func togglePromtOn(s *state.State) {
	if s.Selected.File == nil {
		return
	}
	s.IsPromting = true
	s.Message = "Type 'y' to delete '" + s.Selected.File.Name() + "'"
	output.UpdateScreen(s)
}

func checkPromt(ch rune, s *state.State) {
	s.IsPromting = false
	s.Message = ""

	if ch == 'y' {
		deleteFile(s)
	} else {
		output.UpdateScreen(s)
	}
}

func moveCursorDown(s *state.State) {
	s.Selected.Index = getCursorIndex(s.Selected.Index, len(s.Files), "up")
	s.Selected.File = s.Files[s.Selected.Index]
	output.UpdateScreen(s)
}

func moveCursorUp(s *state.State) {
	s.Selected.Index = getCursorIndex(s.Selected.Index, len(s.Files), "down")
	s.Selected.File = s.Files[s.Selected.Index]
	output.UpdateScreen(s)
}

func getCursorIndex(index int, length int, direction string) int {
	if direction == "up" {
		if index+1 >= length {
			return 0
		}
		return index + 1
	}

	if index == 0 {
		return length - 1
	}
	return index - 1
}

func openFile(file os.FileInfo, s *state.State) {
	if file == nil {
		return
	}
	if file.IsDir() {
		s.Prev.Dir = s.Dir
		s.Prev.File = file
		files.UpdateDir(filepath.Join(s.Dir, file.Name()), s, true)
		output.UpdateScreen(s)
	} else {
		cmd := exec.Command("xdg-open", filepath.Join(s.Dir, file.Name()))
		cmd.Start()
	}
}

func changeDirUp(s *state.State) {
	newDir := path.Dir(s.Dir)

	resetSelected := true

	if newDir == s.Prev.Dir {
		s.Selected.Index = s.Prev.Index
		s.Selected.File = s.Prev.File
		resetSelected = false
	}

	files.UpdateDir(newDir, s, resetSelected)

	output.UpdateScreen(s)
}

func toggleHidden(s *state.State) {
	s.ShowHidden = !s.ShowHidden
	files.UpdateDir(s.Dir, s, false)
	output.UpdateScreen(s)
}

func editFile(s *state.State) {
	editor := os.Getenv("EDITOR")
	cmd := exec.Command(editor, s.Selected.File.Name())
	cmd.Dir = s.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	output.FixScreen()
	output.UpdateScreen(s)
}

func deleteFile(s *state.State) {
	if s.Selected.File == nil {
		return
	}

	err := files.RemoveFile(s.Dir, s.Selected.File)
	if err != nil {
		s.Message = err.Error()
	} else {
		if s.Selected.Index != 0 {
			s.Selected.Index = s.Selected.Index - 1
			s.Selected.File = s.Files[s.Selected.Index]
		} else {
			s.Selected.File = nil
		}
	}

	files.UpdateDir(s.Dir, s, false)
	output.UpdateScreen(s)
}
