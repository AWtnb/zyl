package launchentry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/yaml.v2"
)

func leaf(s string) string {
	ss := strings.Split(filepath.ToSlash(s), "/")
	if 1 < len(ss) {
		return ss[len(ss)-1]
	}
	return s
}

type LaunchEntry struct {
	Path  string
	Alias string
	Depth int
}

func (le *LaunchEntry) resolvePath() {
	le.Path = os.ExpandEnv(le.Path)
}

func (le *LaunchEntry) setAlias() {
	if len(le.Alias) < 1 {
		le.Alias = leaf(le.Path)
	}
}

type LaunchEntries struct {
	entries []LaunchEntry
}

func (les *LaunchEntries) Load(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	entries := []LaunchEntry{}
	if err := yaml.Unmarshal(buf, &entries); err != nil {
		return err
	}
	les.entries = entries
	les.setEditItem(path)
	return nil
}

func (les *LaunchEntries) format() {
	for i := 0; i < len(les.entries); i++ {
		les.entries[i].resolvePath()
		les.entries[i].setAlias()
	}
}

func (les *LaunchEntries) setEditItem(editPath string) {
	ed := []LaunchEntry{{editPath, "EDIT", 0}}
	les.entries = append(ed, les.entries...)
}

func (les LaunchEntries) Select() (le LaunchEntry, err error) {
	les.format()
	candidates := les.entries
	idx, err := fuzzyfinder.Find(candidates, func(i int) string {
		return candidates[i].Alias
	})
	if err != nil {
		return
	}
	le = candidates[idx]
	return
}
