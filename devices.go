package main

import (
	"os/exec"
	"strings"
)

// Device is a struct representing a block device.
type Device struct {
	Name  string
	Model string
	Size  string
}

// GetDevices returns the list of USB devices available to read/write from.
func GetDevices() ([]Device, error) {
	res, err := exec.Command("lsblk", "-o", "KNAME,TYPE,SIZE,MODEL").Output()
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
			device := Device{
				Name: "/dev/" + disk[0],
				Size: disk[2],
			}

			if len(disk) >= 4 && disk[3] != "" {
				device.Model = disk[3]
			}

			disks = append(disks, device)
		}
	}

	return disks, nil
}
