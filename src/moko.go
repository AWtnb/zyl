package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/moko/util"
	"github.com/AWtnb/moko/walk"
	"github.com/go-yaml/yaml"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		datapath string
		filer    string
		all      bool
		exclude  string
	)
	flag.StringVar(&datapath, "datapath", "", "databese yaml file")
	flag.StringVar(&filer, "filer", "explorer.exe", "path of filer")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "search exception (comma-separated)")
	flag.Parse()
	os.Exit(run(datapath, filer, all, exclude))
}

func run(datapath string, filer string, all bool, exclude string) int {
	if !util.IsValidPath(datapath) {
		fmt.Println("cannot find data file...")
		return 1
	}
	if !util.IsValidPath(filer) {
		filer = "explorer.exe"
	}
	lis := loadSource(datapath)
	idx, err := fuzzyfinder.Find(lis, func(i int) string {
		return lis[i].Alias
	})
	if err != nil {
		return 1
	}
	li := lis[idx]
	lp := li.Path
	if util.IsExecutable(lp) {
		util.ExecuteFile(lp)
		return 0
	}
	ld := li.Depth
	if ld < 0 {
		exec.Command(filer, lp).Start()
		return 0
	}
	cs := walk.GetChildItems(lp, ld, all, util.ToSlice(exclude, ","))
	c, err := selectPath(lp, cs)
	if err != nil {
		return 1
	}
	if util.IsExecutable(c) {
		util.ExecuteFile(c)
		return 0
	}
	exec.Command(filer, c).Start()
	return 0
}

func formatChildPath(root string, child string) string {
	rel, _ := filepath.Rel(root, child)
	if fi, _ := os.Stat(child); fi.IsDir() && !util.HasFile(child) {
		return rel + "$"
	}
	return rel
}

func selectPath(root string, paths []string) (string, error) {
	if len(paths) < 1 {
		return root, nil
	}
	if len(paths) == 1 {
		return paths[0], nil
	}
	idx, err := fuzzyfinder.Find(paths, func(i int) string {
		return formatChildPath(root, paths[i])
	})
	if err != nil {
		return "", err
	}
	return paths[idx], nil
}

// loading yaml

type LaunchInfo struct {
	Path  string
	Alias string
	Depth int
}

func loadYaml(fileBuffer []byte) ([]LaunchInfo, error) {
	var data []LaunchInfo
	err := yaml.Unmarshal(fileBuffer, &data)
	if err != nil {
		fmt.Println(err)
	}
	return data, nil
}

func readFile(path string) []LaunchInfo {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return []LaunchInfo{}
	}
	yml, err := loadYaml(buf)
	if err != nil {
		fmt.Println(err)
	}
	return yml
}

func getDisplayName(s string) string {
	if strings.HasPrefix(s, "http") {
		if u, err := url.Parse(s); err == nil {
			return fmt.Sprintf("link[%s/%s]", u.Host, u.RawQuery)
		}
		return s
	}
	return filepath.Base(s)
}

func loadSource(path string) []LaunchInfo {
	var lis []LaunchInfo
	for _, li := range readFile(path) {
		if strings.HasPrefix(li.Path, "http") || util.IsValidPath(li.Path) {
			var l LaunchInfo
			l.Path = li.Path
			l.Depth = li.Depth
			if len(li.Alias) > 0 {
				l.Alias = li.Alias
			} else {
				l.Alias = getDisplayName(li.Path)
			}
			lis = append(lis, l)
		}
	}
	return lis
}
