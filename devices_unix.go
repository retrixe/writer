//go:build !darwin && !windows

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
				Size:  BytesToString(bytes, false),
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
