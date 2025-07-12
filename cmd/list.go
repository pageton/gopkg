package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		tomlPath := core.GetTomlPath(globalFlag)
		cfg, err := core.LoadToml(tomlPath)
		if err != nil {
			fmt.Printf("\033[31m✖️ Failed to load %s: %v\033[0m\n", tomlPath, err)
			return
		}

		lockMap := map[string]core.LockEntry{}
		locks, _ := core.LoadLockFile(globalFlag)
		for _, entry := range locks {
			lockMap[entry.Name] = entry
		}

		modules := make([]string, 0, len(cfg.Dependencies))
		for m := range cfg.Dependencies {
			modules = append(modules, m)
		}
		sort.Strings(modules)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Module", "Declared", "Locked", "Status"})
		table.SetRowLine(true)
		table.SetAutoWrapText(false)

		rows := [][]string{}
		for _, module := range modules {
			declared := cfg.Dependencies[module]
			locked := "—"
			status := "\033[31m✖️ Not installed\033[0m"

			if lock, ok := lockMap[module]; ok {
				locked = lock.Resolved
				cmp := core.CompareVersions(locked, declared)
				switch cmp {
				case 0:
					status = "\033[32m✔️ Up-to-date\033[0m"
				case -1:
					status = "\033[33m⚠️  Outdated\033[0m"
				default:
					status = "\033[34mℹ️  Ahead\033[0m"
				}
			}

			rows = append(rows, []string{module, declared, locked, status})
		}

		table.AppendBulk(rows)
		table.Render()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "List global dependencies")
	rootCmd.AddCommand(listCmd)
}
