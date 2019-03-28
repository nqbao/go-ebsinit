package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type Disk struct {
	Name       string
	Label      string
	UUID       string
	Type       string
	FileSystem string
}

func FindTargetDiskWithContext(ctx context.Context, label string) (Disk, error) {
	for {
		disk := FindTargetDisk(label)

		// we found the disk
		if disk.Name != "" {
			return disk, nil
		}

		// wait for the ebs to be ready
		log.Printf("Waiting for EBS to be ready...")

		select {
		case <-ctx.Done():
			return disk, ctx.Err()
		case <-time.After(5 * time.Second):
			// nothing
		}
	}
}

func FindTargetDisk(label string) Disk {
	disks, _ := ListAllDisks()

	// fist look for the disk with the right label
	if label != "" {
		for _, disk := range disks {
			if disk.Label == label {
				return disk
			}
		}
	}

	// then try to look at the first unformatted disk
	for _, disk := range disks {
		if disk.FileSystem == "data" {
			return disk
		}
	}

	return Disk{}
}

func ListAllDisks() ([]Disk, error) {
	stdout, _, _ := cmd("lsblk", "-d", "-n", "-r", "-o", "NAME")
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	disks := make([]Disk, len(lines))

	for i, name := range lines {
		disks[i] = GetDiskInfo(name)
	}

	return disks, nil
}

func GetDiskInfo(name string) Disk {
	// clear blockcache
	cmd("blkid", "-g")

	deviceName := fmt.Sprintf("/dev/%s", name)

	stdout, _, _ := cmd("blkid", deviceName, "-o", "export")
	disk := Disk{
		Name: name,
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	for _, line := range lines {
		bits := strings.Split(line, "=")

		if len(bits) < 2 {
			continue
		}

		switch bits[0] {
		case "LABEL":
			disk.Label = bits[1]
		case "UUID":
			disk.UUID = bits[1]
		case "TYPE":
			disk.Type = bits[1]
		}
	}

	// get file system on disk
	stdout, _, _ = cmd("file", "-s", deviceName)
	lines = strings.SplitN(strings.TrimSpace(stdout), ": ", 2)
	disk.FileSystem = lines[1]

	return disk
}

func FormatDisk(name string, format string, label string) error {
	log.Printf("Formatting disk %s to %s", name, format)
	deviceName := fmt.Sprintf("/dev/%s", name)
	stdout, stderr, err := cmd("mkfs", "-t", format, deviceName)

	if err != nil {
		log.Printf("Format command error: %v %v", stdout, stderr)
		return err
	}

	log.Printf("Setting device label to %s", label)
	stdout, stderr, err = cmd("e2label", deviceName, label)

	if err != nil {
		log.Printf("Label command error: %v %v", stdout, stderr)
	}

	return err
}
