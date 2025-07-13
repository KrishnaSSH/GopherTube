package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// CheckCommandExists checks if a command exists in PATH
func CheckCommandExists(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("command '%s' not found in PATH", command)
	}
	return nil
}

// ValidateURL checks if a URL is valid
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
