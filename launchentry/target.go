package launchentry

import (
	"os"
	"path/filepath"

	"github.com/AWtnb/go-walk"
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

func (t Target) IsFile() bool {
	fi, err := os.Stat(t.path)
	return err == nil && !fi.IsDir()
}

func (t Target) GetChildItem(all bool, exclude string) (assisted bool, found []string, err error) {
	var d walk.Dir
	d.Init(t.path, all, t.depth, exclude)
	found, err = d.GetChildItemWithEverything()
	assisted = true
	if err != nil || len(found) < 1 {
		assisted = false
		found, err = d.GetChildItem()
	}
	return
}

func (t Target) SelectItem(childPaths []string, prompt string) (string, error) {
	if len(childPaths) == 1 {
		return childPaths[0], nil
	}
	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(t.path, childPaths[i])
		return filepath.ToSlash(rel)
	}, fuzzyfinder.WithPromptString(prompt))
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}
