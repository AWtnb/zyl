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

func (wex WalkException) isSkippable(name string) bool {
	for _, n := range wex.names {
		if n == name {
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

func getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, string(filepath.Separator)), string(filepath.Separator))
}

type ChildItems struct {
	root     string
	maxDepth int
	paths    []string
}

func (ci ChildItems) isSkippableDepth(path string) bool {
	rd := getDepth(ci.root)
	return ci.maxDepth > 0 && getDepth(path)-rd > ci.maxDepth
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
	var wex WalkException
	ci := ChildItems{maxDepth: depth, root: root}
	wex.setNames(exclude, ",")
	var items []string
	// if depth == 0 {
	// 	return items, nil
	// }
	// rd := getDepth(root)
	if !strings.HasPrefix(root, "C:") {
		items = everything.Scan(root, !all)
		if 0 < len(items) {
			ci.paths = wex.filter(items)
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
		if wex.isSkippable(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			items = append(items, path)
		} else {
			if all {
				items = append(items, path)
			}
		}
		return nil
	})
	return items, err
}
