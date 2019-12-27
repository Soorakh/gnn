package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Soorakh/gnn/cursor"
	"github.com/Soorakh/gnn/output"
	"github.com/nsf/termbox-go"
)

type state struct {
	dir               string
	files             []os.FileInfo
	selectedFile      os.FileInfo
	selectedFileIndex int
	showHidden        bool
}

var s *state

func main() {
	dir, _ := os.Getwd()
	s = &state{}
	updateDir(dir, s)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Ch == 'q':
				break loop
			case ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown:
				moveCursorDown(s)
				output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)
			case ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp:
				moveCursorUp(s)
			case ev.Ch == 'l' || ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyEnter:
				changeDir(s.selectedFile, s)
			case ev.Ch == 'h' || ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyBackspace2:
				changeDirUp(s)
			case ev.Ch == '.':
				toggleHidden(s)
			}
		case termbox.EventResize:
			output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)
		}
	}
}

func toggleHidden(s *state) {
	// TODO selected file
	s.showHidden = !s.showHidden
	updateDir(s.dir, s)
	output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)
}

func moveCursorDown(s *state) {
	s.selectedFileIndex = cursor.MoveDown(s.selectedFileIndex, len(s.files))
	s.selectedFile = s.files[s.selectedFileIndex]
}

func moveCursorUp(s *state) {
	s.selectedFileIndex = cursor.MoveUp(s.selectedFileIndex, len(s.files))
	s.selectedFile = s.files[s.selectedFileIndex]
	output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)
}

func getFiles(dir string, showHidden bool) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})

	if !showHidden {
		notHidden := files[:0]
		for _, file := range files {
			if strings.Index(file.Name(), ".") != 0 {
				notHidden = append(notHidden, file)
			}
		}
		return notHidden
	}

	return files
}

func changeDir(file os.FileInfo, s *state) {
	if file == nil {
		return
	}
	if file.IsDir() {
		delimeter := "/"
		if s.dir == "/" {
			delimeter = ""
		}
		updateDir(s.dir+delimeter+file.Name(), s)
		output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)
	} else {
		cmd := exec.Command("xdg-open", s.dir+"/"+file.Name())
		cmd.Start()
	}
}

func changeDirUp(s *state) {
	p := strings.Split(s.dir, "/")
	plen := len(p)
	if plen < 2 {
		return
	}
	p = p[:plen-1]
	if plen > 2 {
		updateDir(strings.Join(p, "/"), s)
	} else {
		updateDir("/", s)
	}

	output.PrintFiles(s.files, s.selectedFile, s.dir, s.selectedFileIndex)
}

func updateDir(d string, s *state) {
	s.dir = d
	s.files = getFiles(d, s.showHidden)
	s.selectedFileIndex = 0
	if len(s.files) > 0 {
		s.selectedFile = s.files[s.selectedFileIndex]
	} else {
		s.selectedFile = nil
	}
}
