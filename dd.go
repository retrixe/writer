package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

// DdProgress is a struct containing progress of the dd operation.
type DdProgress struct {
	Bytes int
	Error error
	Speed string
}

// CopyConvert executes the `dd` Unix utility and provides its output.
//
// Technically, this isn't true anymore, it executes writer itself
// with some special parameters as admin. The new writer process
// wraps `dd` and accepts "stop\n" stdin to terminate dd. This is
// because killing the process doesn't work with pkexec/osascript,
// and this approach enables us to reimplement dd fully.
func CopyConvert(iff string, of string) (chan DdProgress, io.WriteCloser, error) {
	channel := make(chan DdProgress)
	executable, err := os.Executable()
	if err != nil {
		return nil, nil, err
	}
	cmd, err := ElevatedCommand(executable, "dd", iff, of)
	if err != nil {
		return nil, nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	output, input := io.Pipe()
	cmd.Stderr = input
	cmd.Stdout = input
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}
	// Wait for command to exit.
	channelClosed := false
	var mutex sync.Mutex
	go (func() {
		defer input.Close()
		err := cmd.Wait()
		if err != nil {
			channel <- DdProgress{
				Error: err,
			}
		}
		mutex.Lock()
		defer mutex.Unlock()
		channelClosed = true
		close(channel)
	})()
	// Read the output line by line.
	go (func() {
		scanner := bufio.NewScanner(output)
		scanner.Split(ScanCrLines)
		for scanner.Scan() {
			text := scanner.Text()
			println(text)
			firstSpace := strings.Index(text, " ")
			if firstSpace != -1 && strings.HasPrefix(text[firstSpace+1:], "bytes (") {
				// TODO: Probably handle error, but we can't tell full dd behavior without seeing the code.
				bytes, _ := strconv.Atoi(text[:firstSpace])
				split := strings.Split(text, ", ")
				mutex.Lock()
				if channelClosed {
					return // We don't need to unlock as no deadlock is caused here.
				}
				channel <- DdProgress{
					Bytes: bytes,
					Speed: split[len(split)-1],
				}
				mutex.Unlock()
			}
		}
	})()
	return channel, stdin, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanCrLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one carriage return or one mandatory
// newline. In regular expression notation, it is `\r|\n`. The last
// non-empty line of input will be returned even if it has no newline.
func ScanCrLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	} else if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full carriage return-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
