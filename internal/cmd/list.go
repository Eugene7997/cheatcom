package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/Eugene7997/cheatcom/internal/store"
)

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List all saved cheats",
	Args:         cobra.NoArgs,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	s, err := store.Load()
	if err != nil {
		return err
	}

	cheats := s.All()
	if len(cheats) == 0 {
		fmt.Println("No cheats yet. Add one with: chc add \"<command>\" -d \"<description>\"")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, c := range cheats {
		tags := ""
		if len(c.Tags) > 0 {
			tags = "#" + strings.Join(c.Tags, " #")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.ID, c.Command, tags, c.Description)
	}
	return w.Flush()
}
