package walk

import (
	"os"
	"strings"
)

type WalkException struct {
	names []string
}

func (wex *WalkException) SetName(s string) {
	wex.names = append(wex.names, strings.TrimSpace(s))
}

func (wex *WalkException) SetNames(s string, sep string) {
	if len(s) < 1 {
		return
	}
	for _, elem := range strings.Split(s, sep) {
		wex.SetName(elem)
	}
}

func (wex WalkException) Contains(name string) bool {
	for _, n := range wex.names {
		if n == name {
			return true
		}
	}
	return false
}

func (wex WalkException) isSkippablePath(path string) bool {
	sep := string(os.PathSeparator)
	if strings.Contains(path, sep+".") {
		return true
	}
	for _, n := range wex.names {
		if strings.Contains(path, sep+n+sep) || strings.HasSuffix(path, n) {
			return true
		}
	}
	return false
}

func (wex WalkException) Filter(paths []string) []string {
	if len(wex.names) < 1 {
		return paths
	}
	sl := []string{}
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if wex.isSkippablePath(p) {
			continue
		}
		sl = append(sl, p)
	}
	return sl
}
