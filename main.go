package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/sqweek/dialog"
	"github.com/webview/webview"
)

const html = `
<html lang="en">
<head>
  <meta charset="UTF-8">
  <!-- Use minimum-scale=1 to enable GPU rasterization -->
  <meta
    name='viewport'
    content='user-scalable=0, initial-scale=1, minimum-scale=1, width=device-width, height=device-height'
  />
</head>
<body style="margin: 0;"><div id="app"></div><script>initiateReact()</script></body>
</html>
`

var w webview.WebView

// TODO: Do these even need to be bound? Can't JavaScript store them and send
// them back when needed? This is creating 2 unnecessary pieces of state that
// are prone to desync.
var file = ""
var selectedDevice = ""

// ParseToJsString takes a string and escapes slashes and double-quotes,
// and converts it to a string that can be passed to JavaScript.
func ParseToJsString(s string) string {
	return "\"" + strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"") + "\""
}

// SetFile sets the value of the file variable in both Go and React.
func SetFile(value string) {
	file = value
	w.Eval("setFileReact(" + ParseToJsString(value) + ")")
}

// SetSelectedDevice sets the value of the selectedDevice variable in both Go and React.
func SetSelectedDevice(value string) {
	selectedDevice = value
	w.Eval("setSelectedDeviceReact(" + ParseToJsString(value) + ")")
}

func main() {
	debug := true // TODO
	w = webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Writer")
	w.SetSize(420, 210, webview.HintNone)

	// Bind variables.
	w.Bind("setFileGo", func(newFile string) {
		file = newFile
	})
	w.Bind("setSelectedDeviceGo", func(newSelectedDevice string) {
		selectedDevice = newSelectedDevice
	})

	// Bind a function to initiate React via webview.Eval.
	w.Bind("initiateReact", func() {
		w.Eval(js)
		w.SetSize(420, 100, webview.HintMin)
		// Call setDevicesReact.
		devices := GetDevices()
		jsonifiedDevices := make([]string, len(devices))
		for index, device := range devices {
			jsonifiedDevices[index] = ParseToJsString(device)
		}
		w.Eval("setDevicesReact([" + strings.Join(jsonifiedDevices, ", ") + "])")
		SetSelectedDevice(devices[0])
	})

	// Bind a function to prompt for file.
	w.Bind("promptForFile", func() {
		homedir, err := os.UserHomeDir()
		if err != nil {
			homedir = "/"
		}
		filename, err := dialog.File().Title("Select image to flash").SetStartDir(homedir).Filter("Disk image file", "*iso", "*img").Load()
		if err != nil {
			w.Eval("setDialogReact(" + ParseToJsString("Error: "+err.Error()) + ")")
			return
		}
		SetFile(filename) // Send this back to React as well.
	})

	// TODO: Bind privilege escalation.
	// https://github.com/jorangreef/sudo-prompt
	// https://stackoverflow.com/questions/31558066/how-to-ask-for-administer-privileges-on-windows-with-go

	// Bind flashing.
	w.Bind("flash", func() {
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
