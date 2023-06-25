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

func executeFile(path string) {
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func openDir(filer string, path string) {
	exec.Command(filer, path).Start()
}

func isValidPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func toSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}

func selectChildren(root string, paths []string) (string, error) {
	if len(paths) == 1 {
		return paths[0], nil
	}
	idx, err := fuzzyfinder.Find(paths, func(i int) string {
		rel, _ := filepath.Rel(root, paths[i])
		return rel
	})
	if err != nil {
		return "", err
	}
	return paths[idx], nil
}

func run(src string, filer string, all bool, exclude string) int {
	if len(src) < 1 {
		p, _ := os.Executable()
		src = filepath.Join(filepath.Dir(p), "launch.yaml")
	}
	les, err := launchentry.Load(src)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	idx, err := fuzzyfinder.Find(les, func(i int) string {
		return les[i].Alias
	})
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}
	lp := les[idx].Path
	ld := les[idx].Depth
	if strings.HasPrefix(lp, "http") {
		executeFile(lp)
		return 0
	}
	if !isValidPath(lp) {
		err := fmt.Errorf("invalid path: %s", lp)
		fmt.Println(err.Error())
		return 1
	}
	if !isDir(lp) {
		executeFile(lp)
		return 0
	}
	cs, err := walk.GetChildItems(lp, ld, all, toSlice(exclude, ","))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(cs) < 1 {
		openDir(filer, lp)
		return 0
	}
	c, err := selectChildren(lp, cs)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err)
		}
		return 1
	}
	if isDir(c) {
		openDir(filer, c)
	} else {
		executeFile(c)
	}
	return 0
}
