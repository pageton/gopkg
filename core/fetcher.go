package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadModuleZip(module, version string) (string, error) {
	if version == "" {
		return "", fmt.Errorf("version is required")
	}

	safeName := strings.ReplaceAll(module, "/", "_")
	homeDir, _ := os.UserHomeDir()

	baseGopkgDir := filepath.Join(homeDir, ".gopkg")
	cacheDir := filepath.Join(baseGopkgDir, "cache")
	cacheFile := filepath.Join(cacheDir, fmt.Sprintf("%s@%s.zip", safeName, version))

	if _, err := os.Stat(cacheFile); err == nil {
		fmt.Printf("\033[36m📦 Using cached %s@%s\033[0m\n", module, version)
		return cacheFile, nil
	}

	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.zip", module, version)

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache dir: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch zip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch zip: status %d", resp.StatusCode)
	}

	out, err := os.Create(cacheFile)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer out.Close()

	fmt.Printf("\033[34m⬇️ Downloading %s@%s...\033[0m\n", module, version)
	progress := &progressWriter{total: resp.ContentLength}
	if _, err = io.Copy(io.MultiWriter(out, progress), resp.Body); err != nil {
		return "", fmt.Errorf("failed to save zip: %w", err)
	}
	fmt.Printf("\r\033[32m✔️ Downloaded and cached %s@%s\033[0m\n", module, version)
	return cacheFile, nil
}

type progressWriter struct {
	written int64
	total   int64
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n := len(b)
	p.written += int64(n)
	percent := float64(p.written) / float64(p.total) * 100
	fmt.Printf("\r\033[36mProgress: %.1f%%\033[0m", percent)
	return n, nil
}
