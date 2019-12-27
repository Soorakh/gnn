package output

import (
	"os"
	"strconv"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func PrintFiles(files []os.FileInfo, selected os.FileInfo, dir string, selectedFileIndex int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	printWide(0, 0, dir, termbox.ColorDefault, termbox.ColorDefault)
	offset := 2
	visibleFiles := getVisible(files, selectedFileIndex)

	for i, f := range visibleFiles {
		prepend := "  "
		if f == selected {
			prepend = "> "
		}
		fg, bg := getColors(f, selected)
		printWide(0, i+offset, prepend+getFileName(f), fg, bg)
	}
	printStatusBar(selectedFileIndex, len(files), getFileName(selected))
}

func printStatusBar(selectedFileIndex int, flen int, selectedFileName string) {
	w, h := termbox.Size()

	for i := 0; i < w; i++ {
		termbox.SetCell(i, h-1, ' ', termbox.ColorWhite, termbox.ColorWhite)
	}

	printWide(
		0,
		h-1,
		strconv.Itoa(selectedFileIndex+1)+
			"/"+strconv.Itoa(flen)+
			" ["+selectedFileName+"]",
		termbox.ColorBlack,
		termbox.ColorWhite)
	printWide(
		w-3,
		h-1,
		getScrollPosition(flen, selectedFileIndex),
		termbox.ColorBlack,
		termbox.ColorWhite)
	termbox.Flush()
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

func Dump(as []string) {
	for i, s := range as {
		x := 40
		for _, r := range s {
			termbox.SetCell(x, i, r, termbox.ColorDefault, termbox.ColorDefault)
			w := runewidth.RuneWidth(r)
			if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
				w = 1
			}
			x += w
		}
	}
}

func getVisible(files []os.FileInfo, selectedFileIndex int) []os.FileInfo {
	_, h := termbox.Size()
	visibleCount := h - 4
	flen := len(files)
	if visibleCount > flen || visibleCount < 0 {
		visibleCount = flen
	}
	offset := 0
	middle := h / 2
	if selectedFileIndex >= visibleCount-middle && visibleCount < flen {
		offset = selectedFileIndex - visibleCount + middle
	}
	tail := visibleCount + offset
	if tail >= flen {
		tail = flen
		offset = flen - visibleCount
	}
	return files[offset:tail]
}

func getScrollPosition(h int, pos int) string {
	if pos == 0 || h == 0 {
		return "top"
	}
	if pos+1 == h {
		return "bot"
	}
	pc := float64(pos+1) * 100 / float64(h)
	return strconv.Itoa(int(pc)) + "%"
}

func getFileName(file os.FileInfo) string {
	suffix := ""
	if file.IsDir() {
		suffix = "/"
	}

	return file.Name() + suffix
}

func getColors(file os.FileInfo, selected os.FileInfo) (termbox.Attribute, termbox.Attribute) {
	switch {
	case file.IsDir() && file == selected:
		return termbox.ColorBlack, termbox.ColorBlue
	case file.IsDir() && file != selected:
		return termbox.ColorBlue, termbox.ColorDefault
	case !file.IsDir() && file == selected:
		return termbox.ColorBlack, termbox.ColorWhite
	case !file.IsDir() && file != selected:
		return termbox.ColorWhite, termbox.ColorDefault
	}
	return termbox.ColorDefault, termbox.ColorDefault
}

func FixScreen() {
	termbox.Sync()
}
