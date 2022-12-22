package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/launchentry"
	"github.com/AWtnb/moko/util"
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

func run(src string, filer string, all bool, exclude string) int {
	if len(src) < 1 {
		p, _ := os.Executable()
		src = filepath.Join(filepath.Dir(p), "launch.yaml")
	}
	les := launchentry.Load(src)
	idx, err := fuzzyfinder.Find(les, func(i int) string {
		return les[i].Alias
	})
	if err != nil {
		return 1
	}
	lp := les[idx].Path
	ld := les[idx].Depth
	if !util.IsDir(lp) {
		if !strings.HasPrefix(lp, "http") && !util.IsValidPath(lp) {
			return 1
		}
		util.ExecuteFile(lp)
		return 0
	}
	if ld == 0 {
		exec.Command(filer, lp).Start()
		return 0
	}
	cs := walk.GetChildItems(lp, ld, all, util.ToSlice(exclude, ","))
	if len(cs) < 1 {
		exec.Command(filer, lp).Start()
		return 0
	}
	c, err := selectChildren(lp, cs)
	if err != nil {
		return 1
	}
	if util.IsDir(c) {
		exec.Command(filer, c).Start()
	} else {
		util.ExecuteFile(c)
	}
	return 0
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
