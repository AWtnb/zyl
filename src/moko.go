package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	if !isValidPath(datapath) {
		fmt.Println("cannot find data file...")
		return 1
	}
	if !isValidPath(filer) {
		filer = "explorer.exe"
	}
	lis := loadSource(datapath)
	idx, err := fuzzyfinder.Find(lis, func(i int) string {
		return lis[i].Name
	})
	if err != nil {
		return 1
	}
	li := lis[idx]
	lp := li.Path
	if isExecutable(lp) {
		exeCmd(lp)
		return 0
	}
	ld := li.Depth
	if ld < 0 {
		exec.Command(filer, lp).Start()
		return 0
	}
	cs := getChildItems(lp, ld, all, toSlice(exclude, ","))
	var c string
	if len(cs) > 1 {
		idx, err := fuzzyfinder.Find(cs, func(i int) string {
			return formatChildPath(lp, cs[i])
		})
		if err != nil {
			return 1
		}
		c = cs[idx]
	} else if len(cs) == 1 {
		c = cs[0]
	} else {
		c = lp
	}
	if isExecutable(c) {
		exeCmd(c)
	} else {
		exec.Command(filer, c).Start()
	}
	return 0
}

func formatChildPath(root string, child string) string {
	rel, _ := filepath.Rel(root, child)
	if fi, _ := os.Stat(child); fi.IsDir() && !hasFile(child) {
		return rel + "$"
	}
	return rel
}

func hasFile(path string) bool {
	nf := 0
	items, err := ioutil.ReadDir(path)
	if err == nil {
		for _, item := range items {
			if !item.IsDir() {
				nf++
			}
		}
	}
	return nf > 0
}

func exeCmd(path string) {
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func isValidPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isExecutable(path string) bool {
	if strings.HasPrefix(path, "http") {
		return true
	}
	fi, _ := os.Stat(path)
	return !fi.IsDir()
}

func toSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}

type LaunchInfo struct {
	Path  string
	Name  string
	Depth int
}

func loadYaml(fileBuffer []byte) ([]LaunchInfo, error) {
	data := make([]LaunchInfo, 100)
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
	n := ""
	if strings.HasPrefix(s, "http") {
		u, err := url.Parse(s)
		if err == nil {
			n = fmt.Sprintf("link[%s/%s]", u.Host, u.RawQuery)
		} else {
			n = s
		}
	} else {
		n = filepath.Base(s)
	}
	return n
}

func loadSource(path string) []LaunchInfo {
	var lis []LaunchInfo
	for _, li := range readFile(path) {
		if strings.HasPrefix(li.Path, "http") || isValidPath(li.Path) {
			var l LaunchInfo
			l.Path = li.Path
			l.Depth = li.Depth
			if len(li.Name) > 0 {
				l.Name = li.Name
			} else {
				l.Name = getDisplayName(li.Path)
			}
			lis = append(lis, l)
		}
	}
	return lis
}

func sliceContains(slc []string, str string) bool {
	for _, v := range slc {
		if v == str {
			return true
		}
	}
	return false
}

func getDepth(path string) int {
	return strings.Count(path, string(filepath.Separator))
}

func getChildItems(root string, depth int, all bool, exclude []string) []string {
	items := make([]string, 0, 10)
	rd := getDepth(root)
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if depth > 0 && getDepth(path)-rd > depth {
			return filepath.SkipDir
		}
		if sliceContains(exclude, info.Name()) {
			return filepath.SkipDir
		}
		if all {
			items = append(items, path)
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			items = append(items, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("failed to traverse directory...")
	}
	return items
}