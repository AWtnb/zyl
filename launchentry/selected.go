package launchentry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Selected struct {
	path  string
	depth int
}

func (se *Selected) SetEntry(entry LaunchEntry) {
	se.path = entry.Path
	se.depth = entry.Depth
}

func (se Selected) Path() string {
	return se.path
}

func (se Selected) IsValid() bool {
	_, err := os.Stat(se.path)
	return err == nil
}

func (se Selected) IsUri() bool {
	return strings.HasPrefix(se.path, "http")
}

func (se Selected) IsFile() bool {
	fi, err := os.Stat(se.path)
	return err == nil && !fi.IsDir()
}

func (se Selected) GetChildItem(all bool, exclude string) (found []string, err error) {
	dw := walk.DirWalker{All: all, Root: se.path}
	dw.SetWalkDepth(se.depth)
	dw.SetWalkException(exclude)
	if strings.HasPrefix(se.path, "C:") {
		return dw.GetChildItem()
	}
	found, err = dw.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		found, err = dw.GetChildItem()
	}
	return
}

func (se Selected) SelectItem(childPaths []string) (string, error) {
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
