package core

import (
	"os"
	"path/filepath"
)

func GetGlobalModulePath(module string) string {
	return filepath.Join(os.Getenv("HOME"), ".gopkg", "modules", module)
}

func GetGlobalModulesPath() string {
	return filepath.Join(os.Getenv("HOME"), ".gopkg", "modules")
}

func GetCacheDir() string {
	return filepath.Join(os.Getenv("HOME"), ".gopkg", "cache")
}

func GetTomlPath(global bool) string {
	if global {
		return filepath.Join(os.Getenv("HOME"), ".gopkg", "gopkg.toml")
	}
	return "gopkg.toml"
}

func GetCurrentDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}

func GetVendorPath() string {
	return "gopkg_modules"
}
