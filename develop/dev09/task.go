package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

type Config struct {
	URL       string
	Recursive bool
	Depth     int
	OutputDir string
}

func parseFlags() Config {
	config := Config{}
	flag.StringVar(&config.URL, "url", "", "URL to download")
	flag.BoolVar(&config.Recursive, "r", false, "Recursive download")
	flag.IntVar(&config.Depth, "l", 1, "Download depth")
	flag.StringVar(&config.OutputDir, "o", ".", "Output directory")
	flag.Parse()

	if config.URL == "" {
		fmt.Println("URL is required")
		flag.Usage()
		os.Exit(1)
	}
	return config
}

func downloadPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func savePage(outputDir, urls string, body []byte) error {
	u, err := url.Parse(urls)
	if err != nil {
		return err
	}

	dir := path.Join(outputDir, u.Host, u.Path)
	if strings.HasSuffix(urls, "/") {
		dir = path.Join(dir, "index.html")
	} else if !strings.Contains(path.Base(dir), ".") {
		dir += ".html"
	}

	err = os.MkdirAll(path.Dir(dir), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(dir, body, 0644)
}

func extractLinks(baseURL string, body []byte) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(string(body)))

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()
			switch t.Data {
			case "a", "link", "script", "img":
				for _, attr := range t.Attr {
					if (t.Data == "a" && attr.Key == "href") ||
						(t.Data == "link" && attr.Key == "href" && attr.Val != "" && (strings.HasSuffix(attr.Val, ".css") || strings.HasPrefix(attr.Val, "http"))) ||
						(t.Data == "script" && attr.Key == "src" && attr.Val != "" && (strings.HasSuffix(attr.Val, ".js") || strings.HasPrefix(attr.Val, "http"))) ||
						(t.Data == "img" && attr.Key == "src") {
						link := attr.Val
						absoluteURL := resolveURL(baseURL, link)
						links = append(links, absoluteURL)
					}
				}
			}
		}
	}
}

func resolveURL(baseURL string, href string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String()
}

func shouldVisit(url string, visited map[string]bool) bool {
	if _, ok := visited[url]; ok {
		return false
	}
	return true
}

func downloadSite(config Config, currentDepth int, visited map[string]bool) {
	if currentDepth > config.Depth {
		return
	}

	body, err := downloadPage(config.URL)
	if err != nil {
		fmt.Println("Error downloading:", err)
		return
	}

	err = savePage(config.OutputDir, config.URL, body)
	if err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	visited[config.URL] = true

	if config.Recursive && currentDepth < config.Depth {
		links, err := extractLinks(config.URL, body)
		if err != nil {
			fmt.Println("Error extracting links:", err)
			return
		}

		for _, link := range links {
			if shouldVisit(link, visited) {
				newConfig := config
				newConfig.URL = link
				downloadSite(newConfig, currentDepth+1, visited)
			}
		}
	}
}

func main() {
	config := parseFlags()
	visited := make(map[string]bool)
	downloadSite(config, 0, visited)
}
