package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/sqweek/dialog"
	"github.com/webview/webview"
)

// TODO: Design UI, check disk vs ISO sizes, allow cancelling flashing and validate image writing.

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

// var file = ""

// ParseToJsString takes a string and escapes slashes and double-quotes,
// and converts it to a string that can be passed to JavaScript.
func ParseToJsString(s string) string {
	return "\"" + strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"") + "\""
}

// SetFile sets the value of the file variable in both Go and React.
// func SetFile(value string) {file = value;w.Eval("setFileReact(" + ParseToJsString(value) + ")")}

func main() {
	debug := true
	w = webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Writer")
	w.SetSize(420, 210, webview.HintMin)

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
			if device.Model == "" {
				jsonifiedDevices[index] = ParseToJsString(device.Name + " (" + device.Size + ")")
			} else {
				jsonifiedDevices[index] = ParseToJsString(device.Name + " (" + device.Model + ", " + device.Size + ")")
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
	w.Bind("flash", func(file string, selectedDevice string) {
		channel, err := CopyConvert(file, selectedDevice)
		if err != nil {
			w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			return
		}
		go (func() {
			for {
				progress, ok := <-channel
				if ok {
					w.Dispatch(func() {
						if progress.Error != nil { // Error is always the last emitted.
							w.Eval("setDialogReact(" + ParseToJsString("Error: "+progress.Error.Error()) + ")")
						} else {
							w.Eval("setProgressReact(" + strconv.Itoa(progress.Bytes) + ")")
						}
					})
				} else {
					break
				}
			}
			w.Dispatch(func() {
				w.Eval("setProgressReact(\"Done!\")")
			})
		})()
	})

	w.Navigate("data:text/html," + html)
	w.Run()
}
