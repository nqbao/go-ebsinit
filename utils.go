package main

import (
	"bytes"
	"os/exec"
)

func cmd(name string, arg ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	c := exec.Command(name, arg...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()

	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}
