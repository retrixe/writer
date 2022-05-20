package main

import (
	"errors"
	"io/fs"
	"os"
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
	res, err := exec.Command("diskutil", "info", "-all").Output()
	if err != nil {
		return nil, err
	}

	availableDisks := strings.Split(string(res), "\n**********\n")
	availableDisks = availableDisks[:len(availableDisks)-1]

	disks := []Device{}

	for _, availableDisk := range availableDisks {
		disk := make(map[string]string)
		lines := strings.Split(availableDisk, "\n")
		for _, rawLine := range lines {
			line := strings.SplitN(strings.TrimSpace(rawLine), ":", 2)
			if len(line) == 2 {
				disk[strings.TrimSpace(line[0])] = strings.TrimSpace(line[1])
			} else {
				disk[strings.TrimSpace(line[0])] = ""
			}
		}
		if disk["Virtual"] != "No" {
			continue
		} else if disk["Whole"] != "Yes" {
			continue
		} else if disk["Device Location"] == "Internal" {
			continue
		}
		splitDiskSize := strings.Split(disk["Disk Size"], " ")
		bytes, _ := strconv.Atoi(splitDiskSize[2][1:])
		device := Device{
			Name:  disk["Device Node"],
			Size:  splitDiskSize[0] + " " + splitDiskSize[1],
			Bytes: bytes,
			Model: disk["Device / Media Name"],
		}
		disks = append(disks, device)
	}

	return disks, nil
}

// UnmountDevice unmounts a block device's partitons before flashing to it.
func UnmountDevice(device string) error {
	// Check if device is mounted.
	stat, err := os.Stat(device)
	if err != nil {
		return err
	} else if stat.Mode().Type()&fs.ModeDevice == 0 {
		return errors.New("provided device is not a block device!")
	}
	// TODO: Discover device partitions and recursively unmount them.
	return nil
}
