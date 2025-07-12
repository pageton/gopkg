package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var addCmd = &cobra.Command{
	Use:   "add <module>@<version>",
	Short: "Add a dependency to gopkg.toml",
	Example: `
  gopkg add github.com/mattn/go-sqlite3@v1.14.17
  gopkg add -g github.com/user/module@latest
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		parts := strings.Split(arg, "@")
		if len(parts) != 2 {
			fmt.Println("\033[31m✖️ Invalid format. Use: gopkg add <module>@<version>\033[0m")
			fmt.Println("\033[34mℹ️  Example: gopkg add github.com/mattn/go-sqlite3@v1.14.17\033[0m")
			return
		}

		module := parts[0]
		version := parts[1]

		tomlPath := core.GetTomlPath(globalFlag)
		cfg, err := core.LoadToml(tomlPath)
		if err != nil {
			fmt.Printf("\033[33m⚠️  %s not found. Creating...\033[0m\n", tomlPath)
			cfg = &core.GopkgToml{
				Name:         "unnamed",
				Dependencies: map[string]string{},
			}
		}

		if cfg.Dependencies == nil {
			cfg.Dependencies = make(map[string]string)
		}
		cfg.Dependencies[module] = version

		if err := core.SaveToml(tomlPath, cfg); err != nil {
			fmt.Printf("\033[31m✖️ Failed to save gopkg.toml: %v\033[0m\n", err)
			return
		}

		fmt.Printf("\033[32m✔️ Added %s@%s to %s\033[0m\n", module, version, tomlPath)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Add dependency to global gopkg.toml")
	rootCmd.AddCommand(addCmd)
}
