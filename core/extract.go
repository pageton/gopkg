package core

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractZip(zipPath, vendorDir, version string, quiet, global bool) error {
	var destRoot string
	if global {
		destRoot = filepath.Join(os.Getenv("HOME"), ".gopkg", "modules", "github.com")
	} else {
		destRoot = filepath.Join("gopkg_modules", "github.com")
	}

	if !quiet {
		fmt.Printf("\033[34mℹ️ Extracting to %s...\033[0m\n", destRoot)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return fmt.Errorf("archive is empty")
	}
	rootDir := strings.SplitN(r.File[0].Name, "/", 2)[0]
	if rootDir == "" {
		return fmt.Errorf("could not determine archive root directory")
	}

	if err := os.MkdirAll(destRoot, 0755); err != nil {
		return err
	}

	extractedFiles := 0
	for _, f := range r.File {
		if !strings.HasPrefix(f.Name, rootDir+"/") {
			continue
		}
		relPath := strings.TrimPrefix(f.Name, rootDir+"/")
		if relPath == "" {
			continue
		}

		destPath := filepath.Join(destRoot, relPath)

		if !strings.HasPrefix(destPath, filepath.Clean(destRoot)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		src, err := f.Open()
		if err != nil {
			return err
		}
		dst, err := os.Create(destPath)
		if err != nil {
			src.Close()
			return err
		}

		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()

		if err != nil {
			return err
		}

		extractedFiles++
	}

	if version != "" {
		oldPath := "./" + vendorDir + "@" + version
		newPath := "./" + vendorDir

		_ = os.RemoveAll(newPath)
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("rename failed: %w", err)
		}
	}

	if !quiet {
		fmt.Printf("\033[32m✔️ Extracted %d files to %s\033[0m\n", extractedFiles, destRoot)
	}
	return nil
}
