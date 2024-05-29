package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestBuiltinCommands(t *testing.T) {
	testCd(t)
	testPwd(t)
	testEcho(t)
	testKill(t)
	testPs(t)
}

func testCd(t *testing.T) {
	initialDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get initial directory: %v", err)
	}

	newDir := "/tmp"
	if err := os.Chdir(newDir); err != nil {
		t.Fatalf("Failed to change directory to %s: %v", newDir, err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	if currentDir != newDir && currentDir != "/private/tmp" {
		t.Fatalf("cd failed: expected %v or /private/tmp, got %v", newDir, currentDir)
	}

	if err := os.Chdir(initialDir); err != nil {
		t.Fatalf("Failed to restore initial directory: %v", err)
	}
}

func testPwd(t *testing.T) {
	expectedDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	output := runShellCommand(t, "pwd")
	if strings.TrimSpace(output) != expectedDir {
		t.Fatalf("pwd failed: expected %v, got %v", expectedDir, output)
	}
}

func testEcho(t *testing.T) {
	expectedOutput := "Hello, World!"
	output := runShellCommand(t, fmt.Sprintf("echo %v", expectedOutput))
	if strings.TrimSpace(output) != expectedOutput {
		t.Fatalf("echo failed: expected %v, got %v", expectedOutput, output)
	}
}

func testKill(t *testing.T) {
	cmd := exec.Command("sleep", "60")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start sleep process: %v", err)
	}
	defer cmd.Process.Kill()

	pid := cmd.Process.Pid
	runShellCommand(t, fmt.Sprintf("kill %d", pid))

	err = cmd.Wait()
	if err == nil || !strings.Contains(err.Error(), "signal: killed") {
		t.Fatalf("kill failed: expected process to be killed, got %v", err)
	}
}

func testPs(t *testing.T) {
	output := runShellCommand(t, "ps")
	if !strings.Contains(output, "PID") {
		t.Fatalf("ps failed: expected PID header in output, got %v", output)
	}
}

func runShellCommand(t *testing.T, command string) string {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run shell command %v: %v", command, err)
	}
	return out.String()
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args[3:]
	runExternalCommand(args)
	os.Exit(0)
}
