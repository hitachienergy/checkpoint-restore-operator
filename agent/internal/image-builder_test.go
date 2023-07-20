package internal

import (
	"hitachienergy.com/cr-operator/agent/config"
	"os"
	"testing"
)

func TestImageBuilder(t *testing.T) {
	file, err := os.ReadFile("test.tar")
	if err != nil {
		t.Error(err)
	}
	err = os.WriteFile("test2.tar", file, 0644)
	if err != nil {
		t.Error(err)
	}
	config.TempDir = "."
	err = CreateOCIImage("test2.tar", "test-container", "checkpoint")
	if err != nil {
		t.Errorf("Expected tar, received %v", err)
	}
}
