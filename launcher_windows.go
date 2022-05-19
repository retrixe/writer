//go:build launcher

package main

import (
	"os"
	"os/exec"
	"path/filepath"

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
	cmd := exec.Command(filepath.Join(folder, "writer", "writer.exe"))
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Process.Release()
	if err != nil {
		panic(err)
	}
}
