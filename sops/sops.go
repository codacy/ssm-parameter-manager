package sops

import "os/exec"

func Decrypt(filePath string) ([]byte, error) {
	return exec.Command("sops", "-d", filePath).Output()
}
