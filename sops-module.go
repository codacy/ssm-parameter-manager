package main

import "os/exec"

func decryptWithSops(filePath string) ([]byte, error) {
	return exec.Command("sops", "-d", filePath).Output()
}
