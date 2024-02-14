package walk

import (
	"os"
	"strings"
)

func GetDepth(path string) int {
	s := string(os.PathSeparator)
	return strings.Count(strings.TrimSuffix(path, s), s)
}

type DirMember struct {
	rootDepth int
	MaxDepth  int
}

func (dm *DirMember) SetRoot(path string) {
	dm.rootDepth = GetDepth(path)
}

func (dm DirMember) IsSkippableDepth(path string) bool {
	return 0 < dm.MaxDepth && dm.MaxDepth < GetDepth(path)-dm.rootDepth
}

func (dm DirMember) FilterByDepth(paths []string) (filteredPaths []string) {
	if dm.MaxDepth < 0 {
		filteredPaths = paths
		return
	}
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if dm.IsSkippableDepth(p) {
			continue
		}
		filteredPaths = append(filteredPaths, p)
	}
	return
}
