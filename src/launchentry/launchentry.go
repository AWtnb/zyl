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

func getDisplayName(s string) string {
	if strings.HasPrefix(s, "http") {
		if u, err := url.Parse(s); err == nil {
			return fmt.Sprintf("link[%s/%s]", u.Host, u.RawQuery)
		}
		return s
	}
	return filepath.Base(s)
}

type LaunchEntry struct {
	Path  string
	Alias string
	Depth int
}

func (le *LaunchEntry) setAlias() {
	if le == nil {
		return
	}
	if len(le.Alias) < 1 {
		le.Alias = getDisplayName(le.Path)
	}
}

func (le *LaunchEntry) resolvePath() {
	if le == nil {
		return
	}
	le.Path = os.ExpandEnv(le.Path)
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
