package launchentry

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	path  string
	alias string
	depth int
}

func (le LaunchEntry) Path() string {
	return le.path
}
func (le LaunchEntry) Alias() string {
	return le.alias
}
func (le LaunchEntry) Depth() int {
	return le.depth
}

func (le *LaunchEntry) setAlias() {
	if le == nil {
		return
	}
	if 0 < len(le.Alias()) {
		return
	}
	if strings.HasPrefix(le.Path(), "http") {
		if u, err := url.Parse(le.Path()); err == nil {
			le.alias = fmt.Sprintf("link[%s/%s]", u.Host, u.RawQuery)
			return
		}
		le.alias = le.Path()
	}
	le.alias = filepath.Base(le.Path())
}

func (le *LaunchEntry) resolvePath() {
	if le == nil {
		return
	}
	le.path = os.ExpandEnv(le.Path())
}

func Load(path string) ([]LaunchEntry, error) {
	les := []LaunchEntry{}
	buf, err := readFile(path)
	if err != nil {
		return les, err
	}
	if err := yaml.Unmarshal(buf, &les); err != nil {
		return les, err
	}
	les = append([]LaunchEntry{{path, "EDIT", 0}}, les...)
	for i := 0; i < len(les); i++ {
		les[i].resolvePath()
		les[i].setAlias()
	}
	return les, err
}
