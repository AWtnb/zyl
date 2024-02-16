package walk

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/everything"
)

type Dir struct {
	All        bool
	Root       string
	member     DirMember
	exeception WalkException
}

func (d *Dir) SetWalkDepth(depth int) {
	dm := DirMember{MaxDepth: depth}
	dm.SetRoot(d.Root)
	d.member = dm
}

func (d *Dir) SetWalkException(exclude string) {
	var wex WalkException
	wex.SetNames(exclude, ",")
	wex.SetName("AppData")
	d.exeception = wex
}

func (d Dir) GetChildItemWithEverything() (assisted bool, found []string, err error) {
	assisted = true
	if d.member.MaxDepth == 0 {
		return
	}
	found, err = everything.Scan(d.Root, !d.All)
	if err != nil {
		return
	}
	if 0 < len(found) {
		found = d.member.FilterByDepth(d.exeception.Filter(found))
	}
	return
}

func (d Dir) GetChildItem() (assisted bool, found []string, err error) {
	assisted = false
	if d.member.MaxDepth == 0 {
		return
	}
	err = filepath.WalkDir(d.Root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.member.IsSkippableDepth(path) {
			return filepath.SkipDir
		}
		if d.exeception.Contains(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			found = append(found, path)
		} else {
			if d.All {
				found = append(found, path)
			}
		}
		return nil
	})
	return
}
