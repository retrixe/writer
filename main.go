package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "embed"

	"github.com/sqweek/dialog"
	"github.com/webview/webview"
)

// TODO: Design UI (with live warnings/errors). Validate written image?

const html = `
<html lang="en">
<head>
  <meta charset="UTF-8">
  <!-- Use minimum-scale=1 to enable GPU rasterization -->
  <meta
    name='viewport'
    content='user-scalable=0, initial-scale=1, minimum-scale=1, width=device-width, height=device-height'
  />
	<style>
	body {
		margin: 0;
		font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",
		  Ubuntu,Cantarell,Oxygen-Sans,"Helvetica Neue",Arial,Roboto,sans-serif;
	}
  </style>
</head>
<body><div id="app"></div><script>initiateReact()</script></body>
</html>
`

var w webview.WebView

//go:embed dist/main.js
var js string

// var file = ""

// ParseToJsString takes a string and escapes slashes and double-quotes,
// and converts it to a string that can be passed to JavaScript.
func ParseToJsString(s string) string {
	return "\"" + strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"") + "\""
}

// SetFile sets the value of the file variable in both Go and React.
// func SetFile(value string) {file = value;w.Eval("setFileReact(" + ParseToJsString(value) + ")")}

func main() {
	if len(os.Args) >= 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		println("writer version v1.0.0-alpha.1")
		return
	}
	debug := true // TODO
	w = webview.New(debug)
	defer w.Destroy()
	w.SetSize(420, 210, webview.HintNone)
	w.SetTitle("Writer")

	// Bind variables.
	// w.Bind("setFileGo", func(newFile string) {file = newFile})

	// Bind a function to initiate React via webview.Eval.
	w.Bind("initiateReact", func() { w.Eval(js) })

	// Bind a function to request refresh of devices attached.
	w.Bind("refreshDevices", func() {
		devices, err := GetDevices()
		if err != nil {
			w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			return
		}
		jsonifiedDevices := make([]string, len(devices))
		for index, device := range devices {
			base := strconv.Itoa(device.Bytes) + " " + device.Name
			if device.Model == "" {
				jsonifiedDevices[index] = ParseToJsString(base + " (" + device.Size + ")")
			} else {
				jsonifiedDevices[index] = ParseToJsString(base + " (" + device.Model + ", " + device.Size + ")")
			}
		}
		// Call setDevicesReact.
		w.Eval("setDevicesReact([" + strings.Join(jsonifiedDevices, ", ") + "])")
		w.Eval("setSelectedDeviceReact(" + jsonifiedDevices[0] + ")")
	})

	// Bind a function to prompt for file.
	w.Bind("promptForFile", func() {
		homedir, err := os.UserHomeDir()
		if err != nil {
			homedir = "/"
		}
		filename, err := dialog.File().Title("Select image to flash").SetStartDir(homedir).Filter("Disk image file", "*iso", "*img").Load()
		if err != nil && err.Error() != "Cancelled" {
			w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			return
		} else if err == nil {
			stat, err := os.Stat(filename)
			if err != nil {
				w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			} else if !stat.Mode().IsRegular() {
				w.Eval("setDialogReact(" + ParseToJsString("Error: Select a regular file!") + ")")
			} else { // Send this back to React.
				w.Eval("setFileSizeReact(" + strconv.Itoa(int(stat.Size())) + ")")
				w.Eval("setFileReact(" + ParseToJsString(filename) + ")")
			}
		}
	})

	// Bind flashing.
	var currentDdProcess *exec.Cmd
	var cancelled bool = false
	w.Bind("flash", func(file string, selectedDevice string) {
		stat, err := os.Stat(file)
		if err != nil {
			w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			return
		} else if !stat.Mode().IsRegular() {
			w.Eval("setDialogReact(" + ParseToJsString("Error: Select a regular file!") + ")")
			return
		} else {
			w.Eval("setFileSizeReact(" + strconv.Itoa(int(stat.Size())) + ")")
		}
		channel, dd, err := CopyConvert(file, selectedDevice)
		currentDdProcess = dd
		if err != nil {
			w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			return
		}
		go (func() {
			errored := false
			for {
				progress, ok := <-channel
				if cancelled {
					cancelled = false
					return
				} else if ok {
					w.Dispatch(func() {
						if progress.Error != nil { // Error is always the last emitted.
							errored = true
							w.Eval("setDialogReact(" + ParseToJsString("Error: "+progress.Error.Error()) + ")")
						} else {
							w.Eval("setSpeedReact(" + ParseToJsString(progress.Speed) + ")")
							w.Eval("setProgressReact(" + strconv.Itoa(progress.Bytes) + ")")
						}
					})
				} else {
					break
				}
			}
			if !errored {
				w.Dispatch(func() { w.Eval("setProgressReact(\"Done!\")") })
			}
		})()
	})

	w.Bind("cancelFlash", func() {
		cancelled = true
		currentDdProcess.Process.Kill()
		w.Dispatch(func() { w.Eval("setProgressReact(\"Cancelled the operation!\")") })
	})

	w.Navigate("data:text/html," + html)
	w.Run()
}
