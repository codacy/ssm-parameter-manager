package sops

import (
	"os"
	"os/exec"
	"testing"
)

// To keep things simple, since we do not have access to KMS keys here
// This code just tries to open a plain text file through sops, validating if sops is instaled and working as expected
func TestFailDecryptPlainFile(t *testing.T) {
	filePath := "/tmp/test-plain-parameter-file.yaml"
	fContent := []byte("{test: {/test/1: test}}")

	err := os.WriteFile(filePath, fContent, 0644)

	if err != nil {
		t.Errorf("Error creating test file.")
	}

	_, err = Decrypt(filePath)

	if exitErr, isExitError := err.(*exec.ExitError); !isExitError || string(exitErr.Stderr) != "sops metadata not found\n" {
		t.Errorf("File was not suposed to have any meta information for sops to use.")
	}
}
