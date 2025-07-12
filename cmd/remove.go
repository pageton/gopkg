package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var removeCmd = &cobra.Command{
	Use:   "remove <module>",
	Short: "Remove a dependency from gopkg.toml",
	Example: `
  gopkg remove github.com/mattn/go-sqlite3
  gopkg remove -g github.com/user/module@latest
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		module := args[0]
		tomlPath := core.GetTomlPath(globalFlag)

		cfg, err := core.LoadToml(tomlPath)
		if err != nil {
			fmt.Printf("\033[31m✖️ Failed to load %s: %v\033[0m\n", tomlPath, err)
			return
		}

		if _, exists := cfg.Dependencies[module]; !exists {
			fmt.Printf("\033[33m⚠️  %s not found in %s\033[0m\n", module, tomlPath)
			return
		}

		delete(cfg.Dependencies, module)
		if err := core.SaveToml(tomlPath, cfg); err != nil {
			fmt.Printf("\033[31m✖️ Failed to save gopkg.toml: %v\033[0m\n", err)
			return
		}
		fmt.Printf("\033[32m✔️ Removed %s from %s\033[0m\n", module, tomlPath)

		entries, _ := core.LoadLockFile(globalFlag)
		var updated []core.LockEntry
		for _, e := range entries {
			if e.Name != module {
				updated = append(updated, e)
			}
		}
		if err := core.WriteLockFile(updated, globalFlag); err == nil {
			fmt.Println("\033[34mℹ️  Updated gopkg.lock\033[0m")
		}

		_ = exec.Command("go", "mod", "edit", "-droprequire="+module).Run()
		_ = exec.Command("go", "mod", "edit", "-dropreplace="+module).Run()
		fmt.Printf("\033[32m✔️ Removed %s from go.mod\033[0m\n", module)
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Remove from global gopkg.toml")
	rootCmd.AddCommand(removeCmd)
}
