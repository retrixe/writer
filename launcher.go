//go:build windows && launcher

package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	_ "embed"
)

//go:embed writer.exe
var writerExe []byte

//go:embed webview.dll
var webviewDll []byte

//go:embed WebView2Loader.dll
var webview2LoaderDll []byte

func main() {
	// Extract writer.exe to %LocalAppData%.
	folder, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(folder, "writer"), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(folder, "writer", "writer.exe"), writerExe, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(folder, "writer", "webview.dll"), webviewDll, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Join(folder, "writer", "WebView2Loader.dll"), webview2LoaderDll, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// https://stackoverflow.com/a/59147866
	verb := "runas"
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	args := strings.Join(os.Args[1:], " ")

	verbPtr, err := syscall.UTF16PtrFromString(verb)
	if err != nil {
		panic(err)
	}
	exePtr, err := syscall.UTF16PtrFromString(exe)
	if err != nil {
		panic(err)
	}
	cwdPtr, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		panic(err)
	}
	argPtr, err := syscall.UTF16PtrFromString(args)
	if err != nil {
		panic(err)
	}

	var showCmd int32 = 0 // SW_NORMAL

	err = ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		panic(err)
	}
}

// The following code is extracted from golang.org/x/sys/windows

func ShellExecute(hwnd uintptr, verb *uint16, file *uint16, args *uint16, cwd *uint16, showCmd int32) error {
	r1, _, e1 := syscall.Syscall6(
		syscall.NewLazyDLL("C:\\Windows\\System32\\shell32.dll").NewProc("ShellExecuteW").Addr(),
		6,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(cwd)),
		uintptr(showCmd),
	)
	if r1 <= 32 {
		return errnoErr(e1)
	}
	return nil
}

// errnoErr returns common boxed Errno values, to prevent allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return syscall.EINVAL
	case 997:
		return syscall.Errno(997)
	}
	return e
}
