package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var (
	cleanCache bool
	cleanLock  bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean installed modules, lockfiles, and cache",
	Run: func(cmd *cobra.Command, args []string) {
		if globalFlag {
			modulesDir := core.GetGlobalModulesPath()
			if err := os.RemoveAll(modulesDir); err == nil {
				fmt.Printf("\033[32m✔️ Removed global modules: %s\033[0m\n", modulesDir)
			} else {
				fmt.Printf("\033[33m⚠️ Failed to remove global modules: %v\033[0m\n", err)
			}
		} else {
			modulesDir := core.GetVendorPath()
			if err := os.RemoveAll(modulesDir); err == nil {
				fmt.Printf("\033[32m✔️ Removed local modules: %s\033[0m\n", modulesDir)
			} else {
				fmt.Printf("\033[33m⚠️ Failed to remove local modules: %v\033[0m\n", err)
			}
		}

		if cleanLock {
			lockPath := core.GetLockFilePath(globalFlag)
			if err := os.Remove(lockPath); err == nil {
				fmt.Printf("\033[32m✔️ Removed lockfile: %s\033[0m\n", lockPath)
			} else {
				fmt.Printf("\033[33m⚠️ No lockfile found at %s\033[0m\n", lockPath)
			}
		} else {
			fmt.Println("\033[34mℹ️ Skipping gopkg.lock (use --lock to remove it)\033[0m")
		}

		if cleanCache {
			cacheDir := core.GetCacheDir()
			if err := os.RemoveAll(cacheDir); err == nil {
				fmt.Printf("\033[32m✔️ Removed cache: %s\033[0m\n", cacheDir)
			} else {
				fmt.Printf("\033[33m⚠️ Failed to remove cache: %v\033[0m\n", err)
			}
		}
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Clean global modules")
	cleanCmd.Flags().BoolVar(&cleanCache, "cache", false, "Also remove module cache")
	cleanCmd.Flags().BoolVar(&cleanLock, "lock", false, "Also remove gopkg.lock")
	rootCmd.AddCommand(cleanCmd)
}
