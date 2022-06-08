package sops

import "os/exec"

// Decrypt uses SOPS through the command line to decrypt a sops encrypted file
func Decrypt(filePath string) ([]byte, error) {
	return exec.Command("sops", "-d", filePath).Output()
}
