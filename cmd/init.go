package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gopkg in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		tomlPath := "gopkg.toml"
		if _, err := os.Stat(tomlPath); err == nil {
			fmt.Println("\033[31m✖️ gopkg.toml already exists\033[0m")
			return
		}

		projectName := filepath.Base(core.GetCurrentDir())
		data := map[string]any{
			"name":         projectName,
			"dependencies": map[string]any{},
		}

		buf, err := toml.Marshal(data)
		if err != nil {
			fmt.Printf("\033[31m✖️ Failed to encode gopkg.toml: %v\033[0m\n", err)
			return
		}

		if err := os.WriteFile(tomlPath, buf, 0644); err != nil {
			fmt.Printf("\033[31m✖️ Failed to write gopkg.toml: %v\033[0m\n", err)
			return
		}

		modulesDir := core.GetVendorPath()
		if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
			if err := os.Mkdir(modulesDir, 0755); err != nil {
				fmt.Printf("\033[31m✖️ Failed to create gopkg_modules directory: %v\033[0m\n", err)
				return
			}
		}

		fmt.Println("\033[32m✔️ Initialized new gopkg.toml project.\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
