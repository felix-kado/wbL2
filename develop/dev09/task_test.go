package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadPage(t *testing.T) {
	expectedBody := "<html><body>Test</body></html>"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(expectedBody))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	body, err := downloadPage(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(body) != expectedBody {
		t.Fatalf("Expected body %s, got %s", expectedBody, string(body))
	}
}
func TestSavePage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testsavepage")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	url := "http://example.com/test"
	body := []byte("<html><body>Test</body></html>")
	err = savePage(tmpDir, url, body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "example.com", "test.html")
	savedBody, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(savedBody) != string(body) {
		t.Fatalf("Expected body %s, got %s", string(body), string(savedBody))
	}
}

func TestExtractLinks(t *testing.T) {
	baseURL := "http://example.com"
	htmlBody := `
		<html>
		<body>
			<a href="/link1">Link 1</a>
			<a href="http://example.com/link2">Link 2</a>
			<img src="/image1.jpg" />
			<script src="script.js"></script>
		</body>
		</html>`
	expectedLinks := []string{
		"http://example.com/link1",
		"http://example.com/link2",
		"http://example.com/image1.jpg",
		"http://example.com/script.js",
	}

	links, err := extractLinks(baseURL, []byte(htmlBody))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(links) != len(expectedLinks) {
		t.Fatalf("Expected %d links, got %d", len(expectedLinks), len(links))
	}

	for i, link := range links {
		if link != expectedLinks[i] {
			t.Fatalf("Expected link %s, got %s", expectedLinks[i], link)
		}
	}
}

func TestResolveURL(t *testing.T) {
	baseURL := "http://example.com/path/"
	tests := []struct {
		href     string
		expected string
	}{
		{"/link1", "http://example.com/link1"},
		{"link2", "http://example.com/path/link2"},
		{"http://example.com/link3", "http://example.com/link3"},
	}

	for _, tt := range tests {
		resolved := resolveURL(baseURL, tt.href)
		if resolved != tt.expected {
			t.Fatalf("Expected %s, got %s", tt.expected, resolved)
		}
	}
}

func TestShouldVisit(t *testing.T) {
	visited := make(map[string]bool)
	visited["http://example.com/visited"] = true

	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com/visited", false},
		{"http://example.com/notvisited", true},
	}

	for _, tt := range tests {
		result := shouldVisit(tt.url, visited)
		if result != tt.expected {
			t.Fatalf("Expected %v, got %v", tt.expected, result)
		}
	}
}
