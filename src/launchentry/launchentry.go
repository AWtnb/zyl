package launchentry

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/util"
	"gopkg.in/yaml.v2"
)

type LaunchEntry struct {
	Path  string
	Alias string
	Depth int
}

func yamlToEntries(path string) ([]LaunchEntry, error) {
	buf := readFile(path)
	var le []LaunchEntry
	err := yaml.Unmarshal(buf, &le)
	return le, err
}

func readFile(path string) []byte {
	buf, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return buf
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

func Load(path string) []LaunchEntry {
	les := []LaunchEntry{{path, "EDIT", 0}}
	es, err := yamlToEntries(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, le := range es {
		le.Path = util.ResolveEnvPath(le.Path)
		if len(le.Alias) < 1 {
			le.Alias = getDisplayName(le.Path)
		}
		les = append(les, le)
	}
	return les
}
