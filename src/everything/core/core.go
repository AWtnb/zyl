// https://github.com/jof4002/Everything/blob/master/everything_windows_amd64.go

package core

import (
	"syscall"
	"unsafe"
)

const (
	EVERYTHING_OK                     = 0 // no error detected
	EVERYTHING_ERROR_MEMORY           = 1 // out of memory.
	EVERYTHING_ERROR_IPC              = 2 // Everything search client is not running
	EVERYTHING_ERROR_REGISTERCLASSEX  = 3 // unable to register window class.
	EVERYTHING_ERROR_CREATEWINDOW     = 4 // unable to create listening window
	EVERYTHING_ERROR_CREATETHREAD     = 5 // unable to create listening thread
	EVERYTHING_ERROR_INVALIDINDEX     = 6 // invalid index
	EVERYTHING_ERROR_INVALIDCALL      = 7 // invalid call
	EVERYTHING_ERROR_INVALIDREQUEST   = 8 // invalid request data, request data first.
	EVERYTHING_ERROR_INVALIDPARAMETER = 9 // bad parameter.
)

const (
	EVERYTHING_REQUEST_FILE_NAME                           = 0x00000001
	EVERYTHING_REQUEST_PATH                                = 0x00000002
	EVERYTHING_REQUEST_FULL_PATH_AND_FILE_NAME             = 0x00000004
	EVERYTHING_REQUEST_EXTENSION                           = 0x00000008
	EVERYTHING_REQUEST_SIZE                                = 0x00000010
	EVERYTHING_REQUEST_DATE_CREATED                        = 0x00000020
	EVERYTHING_REQUEST_DATE_MODIFIED                       = 0x00000040
	EVERYTHING_REQUEST_DATE_ACCESSED                       = 0x00000080
	EVERYTHING_REQUEST_ATTRIBUTES                          = 0x00000100
	EVERYTHING_REQUEST_FILE_LIST_FILE_NAME                 = 0x00000200
	EVERYTHING_REQUEST_RUN_COUNT                           = 0x00000400
	EVERYTHING_REQUEST_DATE_RUN                            = 0x00000800
	EVERYTHING_REQUEST_DATE_RECENTLY_CHANGED               = 0x00001000
	EVERYTHING_REQUEST_HIGHLIGHTED_FILE_NAME               = 0x00002000
	EVERYTHING_REQUEST_HIGHLIGHTED_PATH                    = 0x00004000
	EVERYTHING_REQUEST_HIGHLIGHTED_FULL_PATH_AND_FILE_NAME = 0x00008000
)

var Everything_SetSearch *syscall.LazyProc
var Everything_SetRequestFlags *syscall.LazyProc
var Everything_Query *syscall.LazyProc
var Everything_IsQueryReply *syscall.LazyProc
var Everything_GetNumResults *syscall.LazyProc
var Everything_IsFolderResult *syscall.LazyProc
var Everything_IsFileResult *syscall.LazyProc
var Everything_GetResultFullPathName *syscall.LazyProc

func init() {
	mod := syscall.NewLazyDLL("Everything64.dll")
	if mod != nil {
		Everything_SetSearch = mod.NewProc("Everything_SetSearchW")
		Everything_SetRequestFlags = mod.NewProc("Everything_SetRequestFlags")
		Everything_Query = mod.NewProc("Everything_QueryW")
		Everything_IsQueryReply = mod.NewProc("Everything_QueryW")
		Everything_GetNumResults = mod.NewProc("Everything_GetNumResults")
		Everything_IsFolderResult = mod.NewProc("Everything_IsFolderResult")
		Everything_IsFileResult = mod.NewProc("Everything_IsFileResult")
		Everything_GetResultFullPathName = mod.NewProc("Everything_GetResultFullPathNameW")
	}
}

// function called for each file or directory visited by Walk.
type WalkFunc func(path string, isFile bool) error

// calling walkFn for each file or directory in queried result
func Walk(root string, skipFile bool, walkFn WalkFunc) error {
	err := SetSearch(root)
	if err != nil {
		return err
	}
	SetRequestFlags(EVERYTHING_REQUEST_FILE_NAME | EVERYTHING_REQUEST_PATH)
	Query(true)
	num := GetNumResults()
	for i := 0; i < num; i++ {
		fullname := GetResultFullPathName(i)
		isFile := IsFileResult(i)
		err := walkFn(fullname, isFile)
		if err != nil {
			return err
		}
	}
	return nil

}

// SetSearch void Everything_SetSearchW(LPCWSTR lpString);
func SetSearch(str string) error {
	if Everything_SetSearch != nil {
		p, err := syscall.UTF16PtrFromString(str)
		if err != nil {
			return err
		}
		Everything_SetSearch.Call(uintptr(unsafe.Pointer(p)))
	}
	return nil
}

// SetRequestFlags void Everything_SetRequestFlags(DWORD dwRequestFlags); // Everything 1.4.1
func SetRequestFlags(flags int) {
	if Everything_SetRequestFlags != nil {
		Everything_SetRequestFlags.Call(uintptr(flags))
	}
}

// Query BOOL Everything_QueryW(BOOL bWait);
func Query(bWait bool) bool {
	if Everything_Query != nil {
		var param int
		if bWait {
			param = 1
		}
		r, _, _ := Everything_Query.Call(uintptr(param))
		return r != 0
	}
	return false
}

// GetNumResults DWORD Everything_GetNumResults(void);
func GetNumResults() int {
	if Everything_GetNumResults != nil {
		r, _, _ := Everything_GetNumResults.Call()
		return int(r)
	}
	return 0
}

// GetResultFullPathName DWORD Everything_GetResultFullPathNameW(DWORD dwIndex,LPWSTR wbuf,DWORD wbuf_size_in_wchars);
func GetResultFullPathName(index int) string {
	if Everything_GetResultFullPathName != nil {
		var pathbuf = make([]uint16, 1024)
		Everything_GetResultFullPathName.Call(uintptr(index), uintptr(unsafe.Pointer(&pathbuf[0])), 1023) // bufsize-1
		return syscall.UTF16ToString(pathbuf)
	}
	return ""
}

// IsFolderResult BOOL Everything_IsFolderResult(DWORD dwIndex);
func IsFolderResult(index int) (ret bool) {
	if Everything_IsFolderResult != nil {
		r, _, _ := Everything_IsFolderResult.Call(uintptr(index))
		ret = r != 0
	}
	return
}

// IsFileResult BOOL Everything_IsFileResult(DWORD dwIndex);
func IsFileResult(index int) bool {
	if Everything_IsFileResult != nil {
		r, _, _ := Everything_IsFileResult.Call(uintptr(index))
		return r != 0
	}
	return false
}
