package launchentry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

type Target struct {
	path  string
	depth int
}

func (t *Target) SetEntry(entry LaunchEntry) {
	t.path = entry.Path
	t.depth = entry.Depth
}

func (t Target) Path() string {
	return t.path
}

func (t Target) IsValid() bool {
	_, err := os.Stat(t.path)
	return err == nil
}

func (t Target) IsUri() bool {
	return strings.HasPrefix(t.path, "http")
}

func (t Target) IsFile() bool {
	fi, err := os.Stat(t.path)
	return err == nil && !fi.IsDir()
}

func (t Target) GetChildItem(all bool, exclude string) (found []string, err error) {
	d := walk.Dir{All: all, Root: t.path}
	d.SetWalkDepth(t.depth)
	d.SetWalkException(exclude)
	if strings.HasPrefix(t.path, "C:") {
		return d.GetChildItem()
	}
	found, err = d.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		found, err = d.GetChildItem()
	}
	return
}

func (t Target) SelectItem(childPaths []string) (string, error) {
	if len(childPaths) == 1 {
		return childPaths[0], nil
	}
	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(t.path, childPaths[i])
		return rel
	})
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}
