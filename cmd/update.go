package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/pageton/gopkg/core"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dependencies to their latest versions",
	Example: `
		gopkg update
		gopkg update -g
	  gopkg update github.com/golang-jwt/jwt/v5@latest
	`,
	Run: func(cmd *cobra.Command, args []string) {
		tomlPath := core.GetTomlPath(globalFlag)
		cfg, err := core.LoadToml(tomlPath)
		if err != nil {
			fmt.Printf("\033[31m✖️ Failed to load %s: %v\033[0m\n", tomlPath, err)
			return
		}

		locks, _ := core.LoadLockFile(globalFlag)
		lockMap := map[string]core.LockEntry{}
		for _, e := range locks {
			lockMap[e.Name] = e
		}

		allModules := make([]string, 0, len(cfg.Dependencies))
		for m := range cfg.Dependencies {
			allModules = append(allModules, m)
		}
		sort.Strings(allModules)

		explicitUpdates := map[string]string{}
		if len(args) > 0 {
			for _, arg := range args {
				if strings.Contains(arg, "@") {
					parts := strings.SplitN(arg, "@", 2)
					name, version := parts[0], parts[1]
					explicitUpdates[name] = version
				} else {
					explicitUpdates[arg] = "latest"
				}
			}
		} else {
			for _, m := range allModules {
				explicitUpdates[m] = "latest"
			}
		}

		fmt.Println("\n🔍 Checking for updates...")
		updated := false
		results := [][]string{}
		toInstall := []string{}

		for _, mod := range allModules {
			current := cfg.Dependencies[mod]
			ver := explicitUpdates[mod]
			if ver == "" {
				continue
			}

			var target string
			if ver == "latest" {
				latest, err := core.ResolveLatestVersion(mod)
				if err != nil {
					results = append(results, []string{mod, current, "—", "\033[31m✖️ Failed\033[0m"})
					continue
				}
				target = latest
			} else {
				target = ver
			}

			cmp := core.CompareVersions(current, target)
			status := "\033[32mUp-to-date\033[0m"
			if cmp < 0 || ver != "latest" {
				status = "\033[33mUpdate available\033[0m"
				cfg.Dependencies[mod] = target
				toInstall = append(toInstall, mod+"@"+target)
				updated = true
			}
			results = append(results, []string{mod, current, target, status})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Module", "Declared", "Latest", "Status"})
		table.SetRowLine(true)
		table.SetAutoWrapText(false)
		table.AppendBulk(results)
		table.Render()

		if updated {
			if err := core.SaveToml(tomlPath, cfg); err != nil {
				fmt.Printf("\033[31m✖️ Failed to save updated gopkg.toml: %v\033[0m\n", err)
				return
			}

			fmt.Println("\n📦 Installing updated modules...")
			for i, mod := range toInstall {
				fmt.Printf("[%d/%d] Installing %s... ", i+1, len(toInstall), mod)
				c := exec.Command(os.Args[0], "install", mod)
				if globalFlag {
					c.Args = append(c.Args, "--global", "-g")
				}
				c.Stdout = nil
				c.Stderr = nil
				err := c.Run()
				if err != nil {
					fmt.Printf("\033[31mFailed\033[0m\n")
				} else {
					fmt.Printf("\033[32mDone\033[0m\n")
				}
				time.Sleep(200 * time.Millisecond)
			}
			fmt.Println("\n✔️ Done.")
		} else {
			fmt.Println("\033[32m✔️ All selected dependencies are up to date.\033[0m")
		}
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "Update global dependencies")
	rootCmd.AddCommand(updateCmd)
}
