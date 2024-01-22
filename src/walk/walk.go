package walk

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/everything"
	"github.com/AWtnb/moko/walk/dirmember"
	"github.com/AWtnb/moko/walk/walkexception"
)

type DirWalker struct {
	All        bool
	Root       string
	member     dirmember.DirMember
	exeception walkexception.WalkException
}

func (dw *DirWalker) ChildItemsHandler(depth int) {
	dm := dirmember.DirMember{MaxDepth: depth, Sep: string(os.PathSeparator)}
	dm.SetRoot(dw.Root)
	dw.member = dm
}

func (dw *DirWalker) ExceptionHandler(exclude string) {
	var wex walkexception.WalkException
	wex.SetNames(exclude, ",")
	dw.exeception = wex
}

func (dw DirWalker) GetChildItemWithEverything() (found []string, err error) {
	if dw.member.MaxDepth == 0 {
		return
	}
	found, err = everything.Scan(dw.Root, !dw.All)
	if err != nil {
		return
	}
	if 0 < len(found) {
		found = dw.member.FilterByDepth(dw.exeception.Filter(found))
	}
	return
}

func (dw DirWalker) GetChildItem() (found []string, err error) {
	if dw.member.MaxDepth == 0 {
		return
	}
	err = filepath.WalkDir(dw.Root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dw.member.IsSkippableDepth(path) {
			return filepath.SkipDir
		}
		if dw.exeception.Contains(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			found = append(found, path)
		} else {
			if dw.All {
				found = append(found, path)
			}
		}
		return nil
	})
	return
}
