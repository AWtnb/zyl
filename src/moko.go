package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		datapath string
		filer    string
		all      bool
		exclude  string
	)
	flag.StringVar(&datapath, "datapath", "", "path to databese text file.\nstyle: filepath(|displayName|depth)\n")
	flag.StringVar(&filer, "filer", "explorer.exe", "path of filer")
	flag.BoolVar(&all, "all", false, "search include file")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
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
	launchItems := loadSource(datapath)
	idx, err := fuzzyfinder.Find(launchItems, func(i int) string {
		return launchItems[i].Name
	})
	if err != nil {
		return 1
	}
	root := launchItems[idx]
	rp := root.Path
	if isExecutable(rp) {
		exeCmd(rp)
		return 0
	}
	md := root.MaxDepth
	if md == 0 {
		exec.Command(filer, rp).Start()
		return 0
	}
	cs := getChildItems(rp, md, all, toSlice(exclude, ","))
	var c string
	if len(cs) > 1 {
		idx, err := fuzzyfinder.Find(cs, func(i int) string {
			return formatChildPath(rp, cs[i])
		})
		if err != nil {
			return 1
		}
		c = cs[idx]
	} else if len(cs) == 1 {
		c = cs[0]
	} else {
		c = rp
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

func getDisplayName(s string) string {
	dn := ""
	if strings.HasPrefix(s, "http") {
		u, err := url.Parse(s)
		if err == nil {
			dn = fmt.Sprintf("link[%s/%s]", u.Host, u.RawQuery)
		} else {
			dn = s
		}
	} else {
		dn = filepath.Base(s)
	}
	return dn
}

func toSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}

type LaunchInfo struct {
	Path     string
	Name     string
	MaxDepth int
}

func parseLine(line string, sep string) LaunchInfo {
	sl := toSlice(line, sep)
	if len(sl) > 2 {
		sl = sl[:3]
	} else if len(sl) == 2 {
		sl = append(sl, "-1")
	} else if len(sl) == 1 {
		sl = append(sl, getDisplayName(sl[0]), "-1")
	}
	if len(sl[1]) < 1 {
		sl[1] = getDisplayName(sl[0])
	}
	md := 0
	if i, err := strconv.Atoi(sl[2]); err == nil {
		md = i
	}
	return LaunchInfo{
		Path:     sl[0],
		Name:     sl[1],
		MaxDepth: md,
	}
}

func readFile(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	lines := make([]string, 0, 100)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return lines
}

type LaunchInfos []LaunchInfo

func loadSource(path string) LaunchInfos {
	var lis LaunchInfos
	for _, s := range readFile(path) {
		if len(strings.TrimSpace(s)) < 1 {
			continue
		}
		li := parseLine(s, "|")
		if strings.HasPrefix(li.Path, "http") || isValidPath(li.Path) {
			lis = append(lis, li)
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

func getChildItems(root string, maxDepth int, all bool, exclude []string) []string {
	items := make([]string, 0, 1000)
	rd := getDepth(root)
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if maxDepth > 0 && getDepth(path)-rd > maxDepth {
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
