package main

import (
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func startTestServer(t *testing.T, response string) (net.Listener, string) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			defer conn.Close()
			_, err = conn.Write([]byte(response))
			if err != nil {
				return
			}
		}
	}()

	return listener, listener.Addr().String()
}

func TestTelnetClient(t *testing.T) {
	expectedResponse := "Hello from test server"
	listener, address := startTestServer(t, expectedResponse)
	defer listener.Close()

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatalf("Failed to split host and port: %v", err)
	}

	config := &Config{
		Timeout: 5 * time.Second,
		Host:    host,
		Port:    port,
	}

	// Подключение к тестовому серверу
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Host, config.Port), config.Timeout)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Чтение из сокета и проверка полученного ответа
	buffer := make([]byte, len(expectedResponse))
	_, err = conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	if string(buffer) != expectedResponse {
		t.Fatalf("Expected %q, but got %q", expectedResponse, string(buffer))
	}
}

func TestMainFunction(t *testing.T) {
	expectedResponse := "Hello from test server"
	listener, address := startTestServer(t, expectedResponse)
	defer listener.Close()

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatalf("Failed to split host and port: %v", err)
	}

	os.Args = []string{"cmd", "--timeout=5s", host, port}

	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	go func() {
		main()
		w.Close()
	}()

	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("Failed to read stdout: %v", err)
		}
		if n == 0 {
			break
		}
		output.Write(buf[:n])
	}

	// Restore stdout
	os.Stdout = stdout

	// Check the output for connection message
	if !strings.Contains(output.String(), "Connected to") {
		t.Fatalf("Expected connection message, but got %s", output.String())
	}

	// Check the output for the server response
	if !strings.Contains(output.String(), expectedResponse) {
		t.Fatalf("Expected server response %q, but got %s", expectedResponse, output.String())
	}
}
