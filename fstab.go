package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func EnsureDiskInFstab(disk Disk, mount, format string) error {
	inFstab, err := CheckDiskInFstab(disk)

	if err != nil {
		return err
	}

	if !inFstab {
		log.Printf("Updating fstab")

		// make sure mount point exists
		cmd("mkdir", "-p", mount)

		return UpdateFstab(disk, mount, format)
	}

	return err
}

func CheckDiskInFstab(disk Disk) (bool, error) {
	f, err := os.Open("/etc/fstab")

	if err != nil {
		return false, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, disk.UUID) {
			return true, nil
		}
	}

	err = scanner.Err()
	if err != nil {
		return false, err
	}

	return false, nil
}

func UpdateFstab(disk Disk, mount, format string) error {
	f, err := os.OpenFile("/etc/fstab", (os.O_WRONLY | os.O_APPEND), os.ModePerm)

	if err != nil {
		return err
	}

	_, err = f.WriteString(
		fmt.Sprintf("\nUUID=%s\t%s\t%s\tdefaults,nofail\t0\t2\n", disk.UUID, mount, format),
	)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func MountAll() {
	cmd("mount", "-a")
}
