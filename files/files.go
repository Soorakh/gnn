package files

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Soorakh/gnn/state"
)

func getFiles(dir string, showHidden bool, search string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// filter by search
	if search != "" {
		filtered := files[:0]
		for _, file := range files {
			if strings.Index(strings.ToLower(file.Name()), strings.ToLower(search)) != -1 {
				filtered = append(filtered, file)
			}
		}
		files = filtered
	}

	// filter by hidden
	if !showHidden {
		notHidden := files[:0]
		for _, file := range files {
			if strings.Index(file.Name(), ".") != 0 {
				notHidden = append(notHidden, file)
			}
		}
		files = notHidden
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})

	return files
}

func UpdateDir(d string, s *state.State, resetSelected bool) {
	ratio := float64(s.Selected.Index) / float64(len(s.Files))

	s.Dir = d
	s.Files = getFiles(d, s.ShowHidden, s.Search.Keyword)
	if !resetSelected && s.Selected.File != nil {
		if s.Selected.Index >= len(s.Files) {
			s.Selected.Index = int(ratio * float64(len(s.Files)))
		}
		for i, v := range s.Files {
			if v.Name() == s.Selected.File.Name() {
				s.Selected.Index = i
				break
			}
		}
	} else {
		s.Selected.Index = 0
	}
	if len(s.Files) > 0 {
		s.Selected.File = s.Files[s.Selected.Index]
	} else {
		s.Selected.File = nil
	}
}

func RemoveFile(path string, file os.FileInfo) error {
	return os.RemoveAll(path + "/" + file.Name())
}

func RenameFile(path string, file os.FileInfo, newname string) error {
	return os.Rename(path+"/"+file.Name(), path+"/"+newname)
}
