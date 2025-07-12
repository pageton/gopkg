package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions <module>",
	Short: "Show all available versions of a module",
	Example: `
		gopkg versions github.com/golang-jwt/jwt/v5
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		module := args[0]
		url := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", module)

		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("\033[31m✖️ Failed to fetch versions for %s\033[0m\n", module)
			return
		}
		defer resp.Body.Close()

		var versions []string
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			versions = append(versions, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("\033[31m✖️ Error reading versions:\033[0m", err)
			return
		}

		sort.Sort(sort.Reverse(sort.StringSlice(versions)))

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"VERSION", "NOTE"})
		table.SetBorder(true)
		table.SetHeaderLine(true)
		table.SetColumnSeparator("|")
		table.SetCenterSeparator("+")
		table.SetRowSeparator("-")
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(false)

		for i, v := range versions {
			note := ""
			if i == 0 {
				note = "\033[32mLatest\033[0m"
			} else if strings.Contains(v, "rc") || strings.Contains(v, "beta") || strings.Contains(v, "alpha") {
				note = "\033[33mPre-release\033[0m"
			} else {
				note = "\033[90mOlder\033[0m"
			}
			table.Append([]string{v, note})
		}

		fmt.Printf("\n\033[34mAvailable versions for %s:\033[0m\n\n", module)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(versionsCmd)
}
