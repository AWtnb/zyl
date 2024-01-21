package launchentry

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/yaml.v2"
)

func readFile(path string) ([]byte, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

type LaunchEntry struct {
	Path  string
	Alias string
	Depth int
}

func (le LaunchEntry) isValid() bool {
	_, err := os.Stat(le.Path)
	return err == nil
}

func (le *LaunchEntry) resolvePath() {
	if le == nil {
		return
	}
	le.Path = os.ExpandEnv(le.Path)
}

func (le *LaunchEntry) setAlias() {
	if le == nil {
		return
	}
	if 0 < len(le.Alias) {
		return
	}
	if strings.HasPrefix(le.Path, "http") {
		if u, err := url.Parse(le.Path); err == nil {
			le.Alias = fmt.Sprintf("link[%s/%s]", u.Host, u.RawQuery)
			return
		}
		le.Alias = le.Path
		return
	}
	le.Alias = filepath.Base(le.Path)
}

type LaunchEntries struct {
	entries []LaunchEntry
}

func (les *LaunchEntries) load(path string) error {
	buf, err := readFile(path)
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

func (les LaunchEntries) validEntries() []LaunchEntry {
	sl := []LaunchEntry{}
	for i := 0; i < len(les.entries); i++ {
		ent := les.entries[i]
		if ent.isValid() {
			sl = append(sl, ent)
		}
	}
	return sl
}

func Select(path string) (LaunchEntry, error) {
	var les LaunchEntries
	if err := les.load(path); err != nil {
		return LaunchEntry{}, err
	}
	les.setEditItem(path)
	les.format()
	candidates := les.validEntries()
	idx, err := fuzzyfinder.Find(candidates, func(i int) string {
		return candidates[i].Alias
	})
	if err != nil {
		return LaunchEntry{}, err
	}
	return candidates[idx], nil
}
