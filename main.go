package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var (
	dir               string
	files             []os.FileInfo
	selectedFile      os.FileInfo
	selectedFileIndex int
	visibleFiles      []os.FileInfo
	showHidden        bool
)

func main() {
	dir, _ = os.Getwd()
	updateDir(dir)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	printFiles(files, selectedFile)

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Ch == 'q':
				break loop
			case ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown:
				moveCursorDown()
			case ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp:
				moveCursorUp()
			case ev.Ch == 'l' || ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyEnter:
				changeDir(selectedFile)
			case ev.Ch == 'h' || ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyBackspace2:
				changeDirUp()
			case ev.Ch == '.':
				toggleHidden()
			}
		case termbox.EventResize:
			printFiles(files, selectedFile)
		}
	}
}

func toggleHidden() {
	// TODO selected file
	showHidden = !showHidden
	updateDir(dir)
	printFiles(files, selectedFile)
}

func moveCursorDown() {
	if selectedFileIndex+1 >= len(files) {
		selectedFileIndex = 0
	} else {
		selectedFileIndex = selectedFileIndex + 1
	}
	selectedFile = files[selectedFileIndex]
	printFiles(files, selectedFile)
}

func moveCursorUp() {
	if selectedFileIndex == 0 {
		selectedFileIndex = len(files) - 1
	} else {
		selectedFileIndex = selectedFileIndex - 1
	}
	selectedFile = files[selectedFileIndex]
	printFiles(files, selectedFile)
}

func getFiles(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

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

func changeDir(file os.FileInfo) {
	if file == nil {
		return
	}
	if file.IsDir() {
		updateDir(dir + "/" + file.Name())
		printFiles(files, selectedFile)
	} else {
		cmd := exec.Command("xdg-open", dir+"/"+file.Name())
		cmd.Start()
	}
}

func changeDirUp() {
	p := strings.Split(dir, "/")
	plen := len(p)
	if plen < 2 {
		return
	}
	p = p[:plen-1]
	if plen > 2 {
		updateDir(strings.Join(p, "/"))
	} else {
		updateDir("/")
	}

	printFiles(files, selectedFile)
}

func updateVisible() {
	_, h := termbox.Size()
	visibleCount := h - 4
	if visibleCount > len(files) || visibleCount < 0 {
		visibleCount = len(files)
	}
	visibleFiles = files[0:visibleCount]
}

func updateDir(d string) {
	dir = d
	files = getFiles(dir)
	selectedFileIndex = 0
	if len(files) > 0 {
		selectedFile = files[selectedFileIndex]
	} else {
		selectedFile = nil
	}
}

func printWide(x, y int, s string, fg termbox.Attribute, bg termbox.Attribute) {
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)
		w := runewidth.RuneWidth(r)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
			w = 1
		}
		x += w
	}
}

func printFiles(files []os.FileInfo, selected os.FileInfo) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	printWide(0, 0, dir, termbox.ColorDefault, termbox.ColorDefault)
	offset := 2
	_, h := termbox.Size()
	updateVisible()

	for i, f := range visibleFiles {
		suffix := ""
		if f.IsDir() {
			suffix = "/"
		}
		if f == selected {
			printWide(0, i+offset, "> "+f.Name()+suffix, termbox.ColorBlack, termbox.ColorBlue)
		} else {
			printWide(0, i+offset, "  "+f.Name()+suffix, termbox.ColorBlue, termbox.ColorDefault)
		}
	}
	printWide(0, h-1, "total: "+strconv.Itoa(len(files)), termbox.ColorBlue, termbox.ColorDefault)
	termbox.Flush()
}
