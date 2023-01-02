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

var envMap = map[string]string{
	"%APPDATA%":     os.Getenv("APPDATA"),
	"%USERNAME%":    os.Getenv("USERNAME"),
	"%USERPROFILE%": os.Getenv("USERPROFILE"),
}

func resolveEnvPath(s string) string {
	for k, v := range envMap {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
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
	lauchEtrs := []LaunchEntry{{path, "EDIT", 0}}
	for _, etr := range rawEtrs {
		etr.Path = resolveEnvPath(etr.Path)
		if len(etr.Alias) < 1 {
			etr.Alias = getDisplayName(etr.Path)
		}
		lauchEtrs = append(lauchEtrs, etr)
	}
	return lauchEtrs, nil
}
