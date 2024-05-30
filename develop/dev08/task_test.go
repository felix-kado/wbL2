package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func captureOutput(f func()) (string, error) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestCDCommand(t *testing.T) {
	cmd := &CDCommand{}
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to return to original directory: %v", err)
		}
	}()

	err := cmd.Execute([]string{"/tmp"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	dir, _ := os.Getwd()
	if dir != "/private/tmp" {
		t.Fatalf("Expected /tmp, got %s", dir)
	}
}

func TestPWDCommand(t *testing.T) {
	cmd := &PWDCommand{}
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to return to original directory: %v", err)
		}
	}()

	err := os.Chdir("/tmp")
	if err != nil {
		t.Fatalf("Could not change to /tmp: %v", err)
	}

	output, err := captureOutput(func() {
		err = cmd.Execute([]string{})
	})

	if err != nil {
		t.Fatalf("Error capturing output: %v", err)
	}

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if strings.TrimSpace(output) != "/private/tmp" {
		t.Fatalf("Expected /tmp, got %s", output)
	}
}

func TestEchoCommand(t *testing.T) {
	cmd := &EchoCommand{}
	output, err := captureOutput(func() {
		err := cmd.Execute([]string{"Hello", "World"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	if err != nil {
		t.Fatalf("Error capturing output: %v", err)
	}

	if strings.TrimSpace(output) != "Hello World" {
		t.Fatalf("Expected 'Hello World', got '%s'", output)
	}
}

func TestKillCommand(t *testing.T) {
	cmd := &KillCommand{}

	// Start a test process to kill
	process := exec.Command("sleep", "10")
	err := process.Start()
	if err != nil {
		t.Fatalf("Could not start process: %v", err)
	}

	err = cmd.Execute([]string{strconv.Itoa(process.Process.Pid)})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = process.Wait()
	if err == nil || !strings.Contains(err.Error(), "signal: killed") {
		t.Fatalf("Expected process to be killed, got %v", err)
	}
}

func TestPSCommand(t *testing.T) {
	cmd := &PSCommand{}
	output, err := captureOutput(func() {
		err := cmd.Execute([]string{})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	if err != nil {
		t.Fatalf("Error capturing output: %v", err)
	}

	if !strings.Contains(output, "PID") {
		t.Fatalf("Expected output to contain 'PID', got '%s'", output)
	}
}

func TestShellExecuteCommand(t *testing.T) {
	shell := NewShell()

	tests := []struct {
		input    string
		expected string
	}{
		{"echo test", "test\n"},
		{"pwd", "/private/tmp\n"},
	}

	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to return to original directory: %v", err)
		}
	}()
	err := os.Chdir("/tmp")
	if err != err {
		t.Fatalf("Failed to change dir to /tmp in test: %v", err)
	}
	for _, test := range tests {
		output, err := captureOutput(func() {
			err := shell.ExecuteCommand(test.input)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})

		if err != nil {
			t.Fatalf("Error capturing output: %v", err)
		}

		if output != test.expected {
			t.Fatalf("Expected '%s', got '%s'", test.expected, output)
		}
	}
}
