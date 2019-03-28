package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	label       string
	format      string
	timeout     int
	mount       string
	showVersion bool
)

var version = "0.1"

func main() {
	flag.StringVar(&label, "label", "DATA", "Target disk label")
	flag.StringVar(&mount, "mount", "/data", "Mount location")
	flag.StringVar(&format, "format", "ext4", "Target disk format")
	flag.IntVar(&timeout, "timeout", 30, "How long in seconds to wait for the EBS to come up")
	flag.BoolVar(&showVersion, "version", false, "Show version")

	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %v\n", version)
		return
	}

	log.Printf("Look for device with label %v with timeout %d seconds", label, timeout)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(int64(timeout))*time.Second)
	disk, err := FindTargetDiskWithContext(ctx, label)

	if err != nil {
		log.Fatal(err)
	} else if disk.Name == "" {
		log.Fatal("Unable to find any suitable device")
	}

	log.Printf("Found device %v, label: %v, type: %v", disk.Name, disk.Label, disk.Type)

	// check if we need to format the disk
	if disk.Type == "" {
		err = FormatDisk(disk.Name, format, label)

		if err != nil {
			log.Fatal(err)
		}
	}

	err = EnsureDiskInFstab(disk, mount, format)
	if err != nil {
		log.Fatalf("Unable to update disk in fstab. Error: %v", err)
	}

	// mount all
	log.Printf("Re-mount all devices")
	MountAll()
}
