package utils

import (
	"os/exec"
)

func RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
