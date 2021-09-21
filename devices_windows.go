package main

import (
	"os/exec"
	"strconv"
	"strings"
)

// Device is a struct representing a block device.
type Device struct {
	Name  string
	Model string
	Size  string
	Bytes int
}

const wmicArgs = "diskdrive get deviceid, mediatype, model, caption, size"

// GetDevices returns the list of USB devices available to read/write from.
func GetDevices() ([]Device, error) {
	res, err := exec.Command("wmic", strings.Split(wmicArgs, " ")...).Output()
	if err != nil {
		return nil, err
	}

	availableDisks := strings.Split(strings.TrimSpace(string(res)), "\n")

	disks := []Device{}

	for _, availableDisk := range availableDisks {
		split := strings.Split(availableDisk, "  ")
		disk := []string{}
		for _, element := range split {
			trimmed := strings.TrimSpace(element)
			if trimmed != "" {
				disk = append(disk, trimmed)
			}
		}
		indexOffset := 0
		if len(disk) == 5 {
			indexOffset = 1
		}
		if disk[1+indexOffset] == "Removable Media" {
			bytes, _ := strconv.Atoi(disk[3+indexOffset])
			device := Device{
				Name:  disk[0+indexOffset],
				Size:  bytesToString(bytes),
				Bytes: bytes,
			}
			if indexOffset == 1 {
				device.Model = disk[0]
			}
			disks = append(disks, device)
		}
	}

	return disks, nil
}

func bytesToString(bytes int) string {
	kb := float64(bytes) / 1000
	mb := kb / 1000
	gb := mb / 1000
	tb := gb / 1000
	if tb >= 1 {
		return strconv.FormatFloat(tb, 'f', 1, 64) + "TB"
	} else if gb >= 1 {
		return strconv.FormatFloat(gb, 'f', 1, 64) + "GB"
	} else if mb >= 1 {
		return strconv.FormatFloat(mb, 'f', 1, 64) + "MB"
	} else if kb >= 1 {
		return strconv.FormatFloat(kb, 'f', 1, 64) + "KB"
	} else {
		return strconv.Itoa(bytes) + "B"
	}
}
