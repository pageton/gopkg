package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type LockEntry struct {
	Name          string `toml:"name"`
	Version       string `toml:"version"`
	Resolved      string `toml:"resolved"`
	Hash          string `toml:"hash"`
	ResolvedTime  string `toml:"resolved_time"`
	InstalledTime string `toml:"installed_time"`
	Source        string `toml:"source"`
}

type LockFile struct {
	Dependencies []LockEntry `toml:"dependencies"`
}

func GetLockFilePath(global bool) string {
	if global {
		return filepath.Join(os.Getenv("HOME"), ".gopkg", "gopkg.lock")
	}
	return "gopkg.lock"
}

func WriteLockFile(entries []LockEntry, global bool) error {
	lock := LockFile{Dependencies: entries}
	lockPath := GetLockFilePath(global)

	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		return fmt.Errorf("failed to create lockfile directory: %w", err)
	}

	f, err := os.Create(lockPath)
	if err != nil {
		return fmt.Errorf("failed to create lockfile: %w", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(lock); err != nil {
		return fmt.Errorf("failed to encode lockfile: %w", err)
	}

	return nil
}

func LoadLockFile(global bool) ([]LockEntry, error) {
	var lock LockFile
	lockPath := GetLockFilePath(global)

	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		return []LockEntry{}, nil
	}

	_, err := toml.DecodeFile(lockPath, &lock)
	if err != nil {
		return nil, fmt.Errorf("failed to decode lockfile: %w", err)
	}

	return lock.Dependencies, nil
}
