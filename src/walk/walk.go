package walk

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/everything"
)

type WalkException struct {
	names []string
}

func (wex *WalkException) setNames(s string, sep string) {
	if len(s) < 1 {
		return
	}
	for _, elem := range strings.Split(s, sep) {
		wex.names = append(wex.names, strings.TrimSpace(elem))
	}
}

func (wex WalkException) contains(name string) bool {
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

func (wex WalkException) filter(paths []string) []string {
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

type ChildItems struct {
	rootDepth int
	maxDepth  int
	sep       string
	paths     []string
}

func (ci *ChildItems) setRoot(path string) {
	ci.rootDepth = ci.getDepth(path)
}

func (ci ChildItems) getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, ci.sep), ci.sep)
}

func (ci ChildItems) isSkippableDepth(path string) bool {
	return 0 < ci.maxDepth && ci.maxDepth < ci.getDepth(path)-ci.rootDepth
}

func (ci ChildItems) filterByDepth() []string {
	if ci.maxDepth == 0 {
		return ci.paths
	}
	sl := []string{}
	for i := 0; i < len(ci.paths); i++ {
		p := ci.paths[i]
		if ci.isSkippableDepth(p) {
			continue
		}
		sl = append(sl, p)
	}
	return sl
}

func GetChildItems(root string, depth int, all bool, exclude string) ([]string, error) {
	ci := ChildItems{maxDepth: depth, sep: string(filepath.Separator)}
	ci.setRoot(root)
	var wex WalkException
	wex.setNames(exclude, ",")
	var found []string
	if !strings.HasPrefix(root, "C:") {
		found = everything.Scan(root, !all)
		if 0 < len(found) {
			ci.paths = wex.filter(found)
			return ci.filterByDepth(), nil
		}
	}
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ci.isSkippableDepth(path) {
			return filepath.SkipDir
		}
		if wex.contains(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			found = append(found, path)
		} else {
			if all {
				found = append(found, path)
			}
		}
		return nil
	})
	return found, err
}
