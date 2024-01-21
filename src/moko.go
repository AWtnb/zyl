package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/launchentry"
	"github.com/AWtnb/moko/walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src     string
		filer   string
		all     bool
		exclude string
	)
	flag.StringVar(&src, "src", "", "source yaml file path")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer path")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "search exception (comma-separated)")
	flag.Parse()
	os.Exit(run(src, filer, all, exclude))
}

type SelectedEntry struct {
	path  string
	depth int
	filer string
}

func (se *SelectedEntry) setEntry(entry launchentry.LaunchEntry) {
	se.path = entry.Path
	se.depth = entry.Depth
}

func (se *SelectedEntry) setFiler(path string) {
	if _, err := os.Stat(path); err == nil {
		se.filer = path
		return
	}
	se.filer = "explorer.exe"
}

func (se SelectedEntry) isValid() bool {
	_, err := os.Stat(se.path)
	return err == nil
}

func (se SelectedEntry) isDir() bool {
	if fi, err := os.Stat(se.path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func (se SelectedEntry) isExecutable() bool {
	return strings.HasPrefix(se.path, "http") || !se.isDir()
}

func (se SelectedEntry) openSelf() {
	if se.isDir() {
		exec.Command(se.filer, se.path).Start()
		return
	}
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", se.path).Start()
}

func (se SelectedEntry) getChildItem(all bool, exclude string) (found []string, err error) {
	de := walk.DirEntry{Root: se.path, All: all, Depth: se.depth, Exclude: exclude}
	if strings.HasPrefix(se.path, "C:") {
		return de.GetChildItem()
	}
	found, err = de.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		found, err = de.GetChildItem()
	}
	return
}

func (se SelectedEntry) selectItem(childPaths []string) (string, error) {
	if len(childPaths) == 1 {
		return childPaths[0], nil
	}
	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(se.path, childPaths[i])
		return rel
	})
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}

func (se SelectedEntry) open(path string) {
	var item SelectedEntry
	item.path = path
	item.setFiler(se.filer)
	item.openSelf()
}

func run(src string, filer string, all bool, exclude string) int {
	if len(src) < 1 {
		p, _ := os.Executable()
		src = filepath.Join(filepath.Dir(p), "launch.yaml")
	}

	var les launchentry.LaunchEntries
	if err := les.Load(src); err != nil {
		fmt.Println(err)
		return 1
	}

	selected, err := les.Select()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}

	var se SelectedEntry
	se.setEntry(selected)
	if !se.isValid() {
		fmt.Printf("invalid path: '%s'\n", selected.Path)
		return 1
	}

	se.setFiler(filer)

	if se.isExecutable() {
		se.openSelf()
		return 0
	}

	cs, err := se.getChildItem(all, exclude)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(cs) < 1 {
		se.openSelf()
		return 0
	}
	c, err := se.selectItem(cs)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}
	se.open(c)
	return 0
}
