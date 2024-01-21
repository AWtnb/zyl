// https://github.com/jof4002/Everything/blob/master/_Example/walk/example.go
package everything

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/moko/everything/core"
)

func Scan(query string, skipFile bool) ([]string, error) {
	sl := []string{}
	if err := checkDll("Everything64.dll"); err != nil {
		return sl, fmt.Errorf("failed to load Everything64.dll")
	}
	core.Walk(query, skipFile, func(path string, isFile bool) error {
		if skipFile && isFile {
			return nil
		}
		sl = append(sl, path)
		return nil
	})
	return sl, nil
}

func getExeDir() string {
	exePath, err := os.Executable()
	if err == nil {
		return filepath.Dir(exePath)
	}
	if exeRealPath, err := filepath.EvalSymlinks(exePath); err == nil {
		return filepath.Dir(exeRealPath)
	}
	return ""
}

func checkDll(name string) error {
	exeDir := getExeDir()
	if len(exeDir) < 1 {
		return fmt.Errorf("failed to detect directory of exe")
	}
	path := filepath.Join(exeDir, name)
	_, err := os.Stat(path)
	return err
}
