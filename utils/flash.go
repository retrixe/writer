package utils

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ErrNotBlockDevice is returned when the specified device is not a block device.
var ErrNotBlockDevice = errors.New("specified device is not a block device")

// RunDd is a wrapper around the `dd` command. This wrapper behaves
// identically to dd, but accepts stdin input "stop\n".
func RunDd(iff string, of string) {
	cmd := exec.Command("dd", "if="+iff, "of="+of, "status=progress", "bs=1M", "conv=fdatasync")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, stdout)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, stderr)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	quit := handleStopInput(func() { cmd.Process.Kill() })
	err = cmd.Wait()
	quit <- true
	if err != nil && cmd.ProcessState.ExitCode() != 0 {
		os.Exit(cmd.ProcessState.ExitCode())
	} else if err != nil {
		panic(err)
	}
}

// FlashFileToBlockDevice is a re-implementation of dd
// in Golang to work cross-platform on Windows as well.
func FlashFileToBlockDevice(iff string, of string) {
	// References to use:
	// https://stackoverflow.com/questions/21032426/low-level-disk-i-o-in-golang
	// https://stackoverflow.com/questions/56512227/how-to-read-and-write-low-level-raw-disk-in-windows-and-go
	quit := handleStopInput(func() { os.Exit(0) })
	filePath, err := filepath.Abs(iff)
	if err != nil {
		log.Fatalln("Unable to resolve path to file.")
	}
	destPath, err := filepath.Abs(of)
	if err != nil {
		log.Fatalln("Unable to resolve path to dest.")
	}
	file, err := os.Open(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Fatalln("This file does not exist!")
	} else if err != nil {
		log.Fatalln("An error occurred while opening the file.", err)
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalln("An error occurred while opening the file.", err)
	} else if fileStat.Mode().IsDir() {
		log.Fatalln("The specified source file is a folder!")
	}
	dest, err := os.OpenFile(destPath, os.O_WRONLY, os.ModePerm) // os.O_RDWR|os.O_EXCL|os.O_CREATE
	if err != nil {
		log.Fatalln("An error occurred while opening the dest.", err)
	}
	defer dest.Close()
	destStat, err := dest.Stat()
	if err != nil {
		log.Fatalln("An error occurred while opening the file.", err)
	} else if destStat.Mode().IsDir() {
		log.Fatalln("The specified destination is a directory!")
	}
	bs := 4096
	if runtime.GOOS == "windows" {
		bs = 512 // TODO: Is this true?
	}
	timer := time.NewTimer(time.Second)
	startTime := time.Now().UnixMilli()
	var total int
	buf := make([]byte, bs)
	for {
		n1, err := file.Read(buf)
		if err != nil {
			if io.EOF == err {
				break
			} else {
				log.Fatalln("Encountered error while reading file!", err)
			}
		}
		n2, err := dest.Write(buf[:n1])
		if err != nil {
			log.Fatalln("Encountered error while writing to dest!", err)
		} else if n2 != n1 {
			log.Fatalln("Read/write mismatch! Is the dest too small!")
		}
		total += n1
		if len(timer.C) > 0 {
			// There's some minor differences in output with dd, mainly decimal places and kB vs KB.
			timeDifference := time.Now().UnixMilli() - startTime
			print(strconv.Itoa(total) + " bytes " +
				"(" + BytesToString(total, false) + ", " + BytesToString(total, true) + ") copied, " +
				strconv.Itoa(int(timeDifference/1000)) + " s, " +
				BytesToString(total/(int(timeDifference)/1000), false) + "/s\r")
			<-timer.C
			timer.Reset(time.Second)
		}
	}
	// t, _ := io.CopyBuffer(dest, file, buf); total = int(t)
	err = dest.Sync()
	if err != nil {
		log.Fatalln("Failed to sync writes to disk!", err)
	} else {
		timeDifference := float64(time.Now().UnixMilli()-startTime) / 1000
		println(strconv.Itoa(total) + " bytes " +
			"(" + BytesToString(total, false) + ", " + BytesToString(total, true) + ") copied, " +
			strconv.FormatFloat(timeDifference, 'f', 3, 64) + " s, " +
			BytesToString(int(float64(total)/timeDifference), false) + "/s")
	}
	quit <- true
}

func handleStopInput(cancel func()) chan bool {
	quit := make(chan bool, 1)
	go (func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-quit:
				return
			default:
				text, err := reader.ReadString('\n')
				if strings.TrimSpace(text) == "stop" {
					cancel()
				} else if err != nil {
					return
				}
			}
		}
	})()
	return quit
}
