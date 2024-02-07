package sys

import "os/exec"

func Open(path string) {
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}
