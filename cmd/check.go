package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if any dependencies have newer versions available",
	Run: func(cmd *cobra.Command, args []string) {
		tomlPath := core.GetTomlPath(globalFlag)
		cfg, err := core.LoadToml(tomlPath)
		if err != nil {
			fmt.Printf("\033[31m✖️ Failed to load %s: %v\033[0m\n", tomlPath, err)
			return
		}

		if len(cfg.Dependencies) == 0 {
			fmt.Printf("\033[34mℹ️  No dependencies found in %s\033[0m\n", tomlPath)
			return
		}

		lockMap := map[string]core.LockEntry{}
		locks, _ := core.LoadLockFile(globalFlag)
		for _, entry := range locks {
			lockMap[entry.Name] = entry
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Module", "Current", "Latest", "Status"})
		table.SetAutoWrapText(false)
		table.SetBorder(true)
		table.SetHeaderLine(true)
		table.SetRowLine(true)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		for module := range cfg.Dependencies {
			locked, found := lockMap[module]
			if !found {
				table.Append([]string{module, "—", "—", "\033[33mNot installed\033[0m"})
				continue
			}

			latest, err := core.ResolveLatestVersion(module)
			if err != nil {
				table.Append([]string{module, locked.Resolved, "—", "\033[31mFailed to fetch\033[0m"})
				continue
			}

			if core.CompareVersions(latest, locked.Resolved) > 0 {
				table.Append([]string{module, locked.Resolved, latest, "\033[33mUpdate available\033[0m"})
			} else {
				table.Append([]string{module, locked.Resolved, latest, "\033[32mUp to date\033[0m"})
			}
		}

		fmt.Println("\n\033[34m📋 Dependency status:\033[0m")
		table.Render()
	},
}

func init() {
	checkCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Check updates in global gopkg.toml")
	rootCmd.AddCommand(checkCmd)
}
