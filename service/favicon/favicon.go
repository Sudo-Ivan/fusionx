package favicon

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Service struct {
	cacheDir string
	client   *http.Client
}

func NewService(cacheDir string) *Service {
	return &Service{
		cacheDir: cacheDir,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Service) GetFaviconPath(feedURL string) (string, error) {
	hostname, err := s.extractHostname(feedURL)
	if err != nil {
		return "", fmt.Errorf("failed to extract hostname: %w", err)
	}

	cacheKey := s.getCacheKey(hostname)
	cachedPath := filepath.Join(s.cacheDir, cacheKey+".png")

	if s.fileExists(cachedPath) {
		return cachedPath, nil
	}

	if err := s.ensureCacheDir(); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return s.fetchAndCacheFavicon(hostname, cachedPath)
}

func (s *Service) extractHostname(feedURL string) (string, error) {
	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		return "", err
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return "", fmt.Errorf("invalid URL: no hostname")
	}

	hostname = strings.TrimPrefix(hostname, "www.")
	
	if strings.Contains(hostname, "rsshub") {
		return s.mapRSSHubToOriginal(parsedURL.Path, hostname)
	}

	return hostname, nil
}

func (s *Service) mapRSSHubToOriginal(path, fallback string) (string, error) {
	rssHubMap := map[string]string{
		"/papers/category/arxiv": "arxiv.org",
		"/trendingpapers/papers": "arxiv.org",
		"/github":                "github.com",
		"/google":                "google.com",
		"/dockerhub":             "hub.docker.com",
		"/imdb":                  "imdb.com",
		"/hackernews":            "news.ycombinator.com",
		"/phoronix":              "phoronix.com",
		"/rsshub":                "rsshub.app",
		"/twitch":                "twitch.tv",
		"/youtube":               "youtube.com",
	}

	for prefix, hostname := range rssHubMap {
		if strings.HasPrefix(path, prefix) {
			return hostname, nil
		}
	}

	return fallback, nil
}

func (s *Service) getCacheKey(hostname string) string {
	var h1, h2 uint32 = 0xdeadbeef, 0x41c6ce57
	for i := 0; i < len(hostname); i++ {
		ch := uint32(hostname[i])
		h1 = h1 ^ ch
		h1 = h1 * 2654435761
		h2 = h2 ^ ch
		h2 = h2 * 1597334677
	}
	h1 = (h1 ^ (h1 >> 16)) * 2246822507 ^ (h2 ^ (h2 >> 13)) * 3266489909
	h2 = (h2 ^ (h2 >> 16)) * 2246822507 ^ (h1 ^ (h1 >> 13)) * 3266489909
	hash := uint64(h2&0x1fffff)<<32 + uint64(h1)
	return fmt.Sprintf("%08x", hash&0xffffffff)[:8]
}

func (s *Service) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (s *Service) ensureCacheDir() error {
	return os.MkdirAll(s.cacheDir, 0750)
}

func (s *Service) fetchAndCacheFavicon(hostname, cachePath string) (string, error) {
	faviconURLs := []string{
		fmt.Sprintf("https://%s/favicon.ico", hostname),
		fmt.Sprintf("https://%s/favicon.png", hostname),
		fmt.Sprintf("https://www.google.com/s2/favicons?sz=32&domain=%s", hostname),
	}

	// First try to find favicons from the website's HTML
	if feedFavicons := s.findFaviconsFromWebsite(hostname); len(feedFavicons) > 0 {
		faviconURLs = append(feedFavicons, faviconURLs...)
	}

	for _, faviconURL := range faviconURLs {
		if err := s.downloadFavicon(faviconURL, cachePath); err == nil {
			return cachePath, nil
		}
	}

	return s.CreateDefaultFavicon(cachePath)
}

func (s *Service) downloadFavicon(faviconURL, cachePath string) error {
	resp, err := s.client.Get(faviconURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// #nosec G304 - cachePath is constructed from sanitized hostname hash
	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *Service) CreateDefaultFavicon(cachePath string) (string, error) {
	// Copy the feed.png file as the default favicon
	sourcePath := "feed.png"
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		// Fallback to a simple inline PNG if feed.png doesn't exist
		defaultFaviconData := []byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
			0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10, 0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x91, 0x68,
			0x36, 0x00, 0x00, 0x00, 0x19, 0x74, 0x45, 0x58, 0x74, 0x53, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72,
			0x65, 0x00, 0x41, 0x64, 0x6f, 0x62, 0x65, 0x20, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x61,
			0x64, 0x79, 0x71, 0xc9, 0x65, 0x3c, 0x00, 0x00, 0x00, 0x32, 0x49, 0x44, 0x41, 0x54, 0x78, 0x44,
			0x62, 0xfc, 0x3f, 0x95, 0x9f, 0x01, 0x37, 0x60, 0x62, 0xc0, 0x0b, 0x46, 0xaa, 0x34, 0x40, 0x80,
			0x01, 0x00, 0x06, 0x50, 0x4e, 0x20, 0x3e, 0x28, 0x84, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
			0x44, 0xae, 0x42, 0x60, 0x82,
		}
		return cachePath, os.WriteFile(cachePath, defaultFaviconData, 0600)
	}
	defer sourceFile.Close()

	// #nosec G304 - cachePath is constructed from sanitized hostname hash
	destFile, err := os.Create(cachePath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return cachePath, err
}

func (s *Service) findFaviconsFromWebsite(hostname string) []string {
	websiteURL := fmt.Sprintf("https://%s", hostname)

	resp, err := s.client.Get(websiteURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// Read only the first 50KB to avoid reading entire large pages
	body := make([]byte, 51200)
	n, err := io.ReadFull(resp.Body, body)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil
	}
	content := string(body[:n])

	// Look for favicon links in HTML head
	faviconRegex := regexp.MustCompile(`(?i)<link[^>]*rel=["'](?:icon|shortcut icon)["'][^>]*href=["']([^"']+)["'][^>]*>`)
	matches := faviconRegex.FindAllStringSubmatch(content, -1)

	var faviconURLs []string
	for _, match := range matches {
		if len(match) > 1 {
			faviconURL := match[1]
			// Convert relative URLs to absolute
			if !strings.HasPrefix(faviconURL, "http") {
				baseURL, err := url.Parse(websiteURL)
				if err != nil {
					continue
				}
				faviconURL = baseURL.ResolveReference(&url.URL{Path: faviconURL}).String()
			}
			faviconURLs = append(faviconURLs, faviconURL)
		}
	}

	return faviconURLs
}
