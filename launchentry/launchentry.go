package launchentry

import (
	"os"
	"path/filepath"
	"sort"
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
	src     string
	entries []LaunchEntry
}

func (les *LaunchEntries) Init(src string) {
	les.src = src
}

func (les *LaunchEntries) Load() error {
	buf, err := os.ReadFile(les.src)
	if err != nil {
		return err
	}
	entries := []LaunchEntry{}
	if err := yaml.Unmarshal(buf, &entries); err != nil {
		return err
	}
	les.entries = entries
	return nil
}

func (les *LaunchEntries) sort() {
	sort.Slice(les.entries, func(i, j int) bool {
		return les.entries[i].Alias < les.entries[j].Alias
	})
}

func (les *LaunchEntries) format() {
	for i := 0; i < len(les.entries); i++ {
		les.entries[i].resolvePath()
		les.entries[i].setAlias()
	}
	les.sort()
	les.setEditItem(les.src)
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
