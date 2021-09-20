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

// GetDevices returns the list of USB devices available to read/write from.
func GetDevices() ([]Device, error) {
	res, err := exec.Command("lsblk", "-b", "-o", "KNAME,TYPE,SIZE,MODEL").Output()
	if err != nil {
		return nil, err
	}

	bootDevices, err := exec.Command("df", "/", "/home").Output()
	if err != nil {
		return nil, err
	}

	splitBoot := strings.Split(string(bootDevices), "\n")
	rootPart := splitBoot[1]
	homePart := splitBoot[2]

	files := strings.Split(string(res), "\n")
	files = files[:len(files)-1]

	disks := []Device{}

	for _, file := range files {
		disk := strings.Fields(file)
		if disk[1] == "disk" && !strings.HasPrefix(disk[0], "zram") &&
			!strings.HasPrefix(rootPart, "/dev/"+disk[0]) &&
			!strings.HasPrefix(homePart, "/dev/"+disk[0]) {
			bytes, _ := strconv.Atoi(disk[2])
			device := Device{
				Name:  "/dev/" + disk[0],
				Size:  bytesToString(bytes),
				Bytes: bytes,
			}

			if len(disk) >= 4 && disk[3] != "" {
				device.Model = disk[3]
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
