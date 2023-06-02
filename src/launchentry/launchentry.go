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

func Load(path string) ([]LaunchEntry, error) {
	rawEtrs := []LaunchEntry{}
	buf, err := readFile(path)
	if err != nil {
		return rawEtrs, err
	}
	if err := yaml.Unmarshal(buf, &rawEtrs); err != nil {
		return rawEtrs, err
	}
	launchEtrs := []LaunchEntry{{path, "EDIT", 0}}
	for _, etr := range rawEtrs {
		etr.Path = os.ExpandEnv(etr.Path)
		if len(etr.Alias) < 1 {
			etr.Alias = getDisplayName(etr.Path)
		}
		launchEtrs = append(launchEtrs, etr)
	}
	return launchEtrs, nil
}
