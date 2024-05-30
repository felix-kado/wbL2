package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("go", "run", "task.go", "--server", "pool.ntp.org")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Program failed with error: %v", err)
	}

	expected := "Current NTP time from server pool.ntp.org is:"
	if !strings.Contains(string(output), expected) {
		t.Fatalf("Unexpected output: %s\nExpected to contain: %s", output, expected)
	}
}

func TestInvalidServer(t *testing.T) {
	// Запускаем нашу программу с неправильным адресом сервера NTP.
	cmd := exec.Command("go", "run", "task.go", "--server", "invalid.server")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("Expected program to fail with an invalid server, but it did not")
	}

	expected := "Error fetching NTP time from server"
	if !strings.Contains(string(output), expected) {
		t.Fatalf("Unexpected output: %s\nExpected to contain: %s", output, expected)
	}
}
