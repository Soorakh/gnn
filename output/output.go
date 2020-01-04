package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Soorakh/gnn/state"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func PrintFiles(files []os.FileInfo, selected os.FileInfo, dir string, selectedFileIndex int, search state.Search, message string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	printWide(0, 0, dir, termbox.ColorDefault, termbox.ColorDefault)
	offset := 2
	visibleFiles := getVisible(files, selectedFileIndex)

	for i, f := range visibleFiles {
		fg, bg := getColors(f, selected)
		printWide(0, i+offset, " "+getFileName(f), fg, bg)
	}
	printStatusBar(selectedFileIndex, len(files), selected, search, message)
}

func printStatusBar(selectedFileIndex int, flen int, selectedFile os.FileInfo, search state.Search, message string) {
	w, h := termbox.Size()

	for i := 0; i < w; i++ {
		termbox.SetCell(i, h-1, ' ', termbox.ColorBlack, termbox.ColorGreen)
	}

	fs := ""

	if message != "" {
		fs = message
	} else if search.IsActive {
		fs = "/" + search.Keyword
	} else {
		if flen == 0 {
			fs = "0/0"
		} else {
			fs = strconv.Itoa(selectedFileIndex+1) +
				"/" + strconv.Itoa(flen) +
				" " + selectedFile.ModTime().Format("2006-01-02 15:04:05") + " " +
				selectedFile.Mode().String() +
				" " + formatSize(selectedFile.Size())
			if !selectedFile.IsDir() && filepath.Ext(selectedFile.Name()) != "" {
				fs = fs + " " + filepath.Ext(selectedFile.Name())
			}
			fs = fs + " [" + getFileName(selectedFile) + "]"
		}
	}

	printWide(
		0,
		h-1,
		fs,
		termbox.ColorBlack,
		termbox.ColorGreen)
	printWide(
		w-3,
		h-1,
		getScrollPosition(flen, selectedFileIndex),
		termbox.ColorBlack,
		termbox.ColorGreen)
	termbox.Flush()
}

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
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
	if file == nil {
		return ""
	}
	suffix := ""
	if file.IsDir() {
		suffix = "/"
	}

	return file.Name() + suffix
}

func getColors(file os.FileInfo, selected os.FileInfo) (termbox.Attribute, termbox.Attribute) {
	switch {
	case file.IsDir() && file == selected:
		return termbox.ColorBlack | termbox.AttrBold, termbox.ColorBlue | termbox.AttrBold
	case file.IsDir() && file != selected:
		return termbox.ColorBlue | termbox.AttrBold, termbox.ColorDefault | termbox.AttrBold
	case !file.IsDir() && file == selected:
		return termbox.ColorBlack, termbox.ColorGreen
	case !file.IsDir() && file != selected:
		return termbox.ColorGreen, termbox.ColorDefault
	}
	return termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault | termbox.AttrBold
}

func FixScreen() {
	termbox.Sync()
}
