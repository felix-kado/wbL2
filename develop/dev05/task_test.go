package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runGrep(args []string, input string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "task.go"}, args...)...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func TestSimpleMatch(t *testing.T) {
	args := []string{"pattern"}
	input := "This is a test.\nAnother line with pattern.\nAnd another one."
	expected := "Another line with pattern.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestIgnoreCase(t *testing.T) {
	args := []string{"-i", "pattern"}
	input := "This is a test.\nAnother line with Pattern.\nAnd another one."
	expected := "Another line with Pattern.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestInvertMatch(t *testing.T) {
	args := []string{"-v", "pattern"}
	input := "This is a test.\nAnother line with pattern.\nAnd another one."
	expected := "This is a test.\nAnd another one.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestFixedString(t *testing.T) {
	args := []string{"-F", "pattern"}
	input := "This is a test.\nAnother line with pattern.\nAnd another one.\npattern"
	expected := "pattern\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestLineNumber(t *testing.T) {
	args := []string{"-n", "pattern"}
	input := "This is a test.\nAnother line with pattern.\nAnd another one."
	expected := "2: Another line with pattern.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestContext(t *testing.T) {
	args := []string{"-C", "1", "pattern"}
	input := "Line one.\nThis is a test.\nAnother line with pattern.\nAnd another one."
	expected := "This is a test.\nAnother line with pattern.\nAnd another one.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestAfterContext(t *testing.T) {
	args := []string{"-A", "1", "pattern"}
	input := "Line one.\nThis is a test.\nAnother line with pattern.\nAnd another one."
	expected := "Another line with pattern.\nAnd another one.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestBeforeContext(t *testing.T) {
	args := []string{"-B", "1", "pattern"}
	input := "Line one.\nThis is a test.\nAnother line with pattern.\nAnd another one."
	expected := "This is a test.\nAnother line with pattern.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestCount(t *testing.T) {
	args := []string{"-c", "2", "pattern"}
	input := "Line one.\nThis is a test.\n1 Another line with pattern.\n2 Another line with pattern.\n3 Another line with pattern.\n4 Another line with pattern.\nAnd another one."
	expected := "1 Another line with pattern.\n2 Another line with pattern.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestAfterBeforeContext(t *testing.T) {
	args := []string{"-B", "1", "-A", "1", "pattern"}
	input := "Line one.\nThis is a test.\nAnother line with pattern.\nAnd another one."
	expected := "This is a test.\nAnother line with pattern.\nAnd another one.\n"

	output, err := runGrep(args, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}
