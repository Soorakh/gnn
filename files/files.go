package files

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func GetFiles(dir string, showHidden bool) []os.FileInfo {
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
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
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
