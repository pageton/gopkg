package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var (
	globalFlag bool
	autoFlag   bool
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install all dependencies from gopkg.toml",
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		tomlPath := core.GetTomlPath(globalFlag)

		var cfg *core.GopkgToml
		var err error

		if autoFlag {
			imports, err := core.ScanImports(".")
			if err != nil {
				fmt.Println("Failed to scan Go files:", err)
				return
			}

			cfg, _ = core.LoadToml(tomlPath)
			if cfg == nil {
				cfg = &core.GopkgToml{
					Name:         filepath.Base(filepath.Dir(tomlPath)),
					Dependencies: map[string]string{},
				}
			}

			for _, imp := range imports {
				if _, ok := cfg.Dependencies[imp]; !ok {
					cfg.Dependencies[imp] = "latest"
					fmt.Println("➕ Auto-added:", imp)
				}
			}

			if err := core.SaveToml(tomlPath, cfg); err != nil {
				fmt.Println("Failed to update gopkg.toml:", err)
				return
			}
		} else {
			cfg, err = core.LoadToml(tomlPath)
			if err != nil {
				cfg = &core.GopkgToml{
					Name:         filepath.Base(filepath.Dir(tomlPath)),
					Dependencies: map[string]string{},
				}
				if err := core.SaveToml(tomlPath, cfg); err != nil {
					fmt.Printf("Failed to create gopkg.toml: %v\n", err)
					return
				}
			}
		}

		// 🔐 Load lockfile
		lockMap := make(map[string]core.LockEntry)
		if lockEntries, err := core.LoadLockFile(globalFlag); err == nil {
			for _, entry := range lockEntries {
				lockMap[entry.Name] = entry
			}
		}

		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			wd, _ := os.Getwd()
			projectName := filepath.Base(wd)
			_ = exec.Command("go", "mod", "init", projectName).Run()
		}

		fmt.Println("\n🔧 Installing dependencies...")

		modules := make([]string, 0, len(cfg.Dependencies))
		for m := range cfg.Dependencies {
			modules = append(modules, m)
		}
		sort.Strings(modules)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Module", "Version", "Resolved", "Status"})
		table.SetAutoWrapText(false)
		table.SetBorder(true)
		table.SetRowLine(true)

		var newLock []core.LockEntry

		for i, module := range modules {
			version := cfg.Dependencies[module]
			fmt.Printf("[%d/%d] Installing %s@%s... \n", i+1, len(modules), module, version)

			var meta *core.ModuleMetadata
			var resolvedVersion string
			status := ""

			if lockEntry, ok := lockMap[module]; ok && lockEntry.Version == version {
				resolvedVersion = lockEntry.Resolved
				meta = &core.ModuleMetadata{
					Version: lockEntry.Resolved,
					Time:    parseTime(lockEntry.ResolvedTime),
					Hash:    lockEntry.Hash,
				}
				status = "Locked"
			} else {
				meta, err = core.FetchModuleMetadata(module, version)
				if err != nil {
					status = "\033[31mFailed\033[0m"
					table.Append([]string{module, version, "—", status})
					continue
				}
				resolvedVersion = meta.Version
				status = "Installed"
			}

			localPath := core.GetVendorPath()
			if globalFlag {
				localPath = core.GetGlobalModulePath(module)
			} else {
				localPath = filepath.Join(core.GetVendorPath(), module)
			}

			modFile := filepath.Join(localPath, "go.mod")
			if _, err := os.Stat(modFile); err != nil {
				zipPath, err := core.DownloadModuleZip(module, resolvedVersion)
				if err != nil {
					status = "Download"
					table.Append([]string{module, version, resolvedVersion, status})
					continue
				}
				err = core.ExtractZip(zipPath, localPath, resolvedVersion, true, globalFlag)
				if err != nil {
					status = "Extract"
					table.Append([]string{module, version, resolvedVersion, status})
					continue
				}
			}

			relPath := localPath
			if !globalFlag {
				relPath, _ = filepath.Rel(".", localPath)
				relPath = "./" + relPath
			}
			err = core.AddReplaceToGoMod(module, relPath, resolvedVersion)
			if err != nil {
				status = "✖️ Replace"
				table.Append([]string{module, version, resolvedVersion, status})
				continue
			}

			newLock = append(newLock, core.LockEntry{
				Name:          module,
				Version:       version,
				Resolved:      resolvedVersion,
				Source:        "github",
				Hash:          meta.Hash,
				ResolvedTime:  meta.Time.Format(time.RFC3339),
				InstalledTime: time.Now().UTC().Format(time.RFC3339),
			})

			table.Append([]string{module, version, resolvedVersion, status})
			time.Sleep(100 * time.Millisecond)
		}

		table.Render()

		if err := core.WriteLockFile(newLock, globalFlag); err == nil {
			fmt.Println("📌 Updated gopkg.lock")
		}
	},
}

func init() {
	installCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Install dependencies globally to ~/.gopkg/modules")
	installCmd.Flags().
		BoolVar(&autoFlag, "auto", false, "Automatically detect imports from Go files and update gopkg.toml")
	rootCmd.AddCommand(installCmd)
}

func parseTime(t string) time.Time {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Now().UTC()
	}
	return parsed
}
