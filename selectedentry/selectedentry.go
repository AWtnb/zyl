package selectedentry

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/launchentry"
	"github.com/AWtnb/moko/walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

type SelectedEntry struct {
	path  string
	depth int
	filer string
}

func (se *SelectedEntry) SetEntry(entry launchentry.LaunchEntry) {
	se.path = entry.Path
	se.depth = entry.Depth
}

func (se *SelectedEntry) SetFiler(path string) {
	if _, err := os.Stat(path); err == nil {
		se.filer = path
		return
	}
	se.filer = "explorer.exe"
}

func (se SelectedEntry) IsValid() bool {
	_, err := os.Stat(se.path)
	return err == nil
}

func (se SelectedEntry) IsDir() bool {
	if fi, err := os.Stat(se.path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func (se SelectedEntry) IsExecutable() bool {
	return strings.HasPrefix(se.path, "http") || !se.IsDir()
}

func (se SelectedEntry) OpenSelf() {
	if se.IsDir() {
		exec.Command(se.filer, se.path).Start()
		return
	}
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", se.path).Start()
}

func (se SelectedEntry) GetChildItem(all bool, exclude string) (found []string, err error) {
	dw := walk.DirWalker{All: all, Root: se.path}
	dw.ChildItemsHandler(se.depth)
	dw.ExceptionHandler(exclude)
	if strings.HasPrefix(se.path, "C:") {
		return dw.GetChildItem()
	}
	found, err = dw.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		found, err = dw.GetChildItem()
	}
	return
}

func (se SelectedEntry) SelectItem(childPaths []string) (string, error) {
	if len(childPaths) == 1 {
		return childPaths[0], nil
	}
	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(se.path, childPaths[i])
		return rel
	})
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}

func (se SelectedEntry) Open(path string) {
	var item SelectedEntry
	item.path = path
	item.SetFiler(se.filer)
	item.OpenSelf()
}
