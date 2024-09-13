package utils

import (
	"errors"
	"os/exec"
	"strings"
)

func GetProjectName() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("Error running go list: " + err.Error())
	}
	return strings.TrimSpace(string(output)), nil
}
