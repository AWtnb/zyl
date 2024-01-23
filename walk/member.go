package walk

import (
	"strings"
)

type DirMember struct {
	rootDepth int
	MaxDepth  int
	Sep       string
}

func (dm *DirMember) SetRoot(path string) {
	dm.rootDepth = dm.getDepth(path)
}

func (dm DirMember) getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, dm.Sep), dm.Sep)
}

func (dm DirMember) IsSkippableDepth(path string) bool {
	return 0 < dm.MaxDepth && dm.MaxDepth < dm.getDepth(path)-dm.rootDepth
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
