package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Eugene7997/cheatcom/internal/store"
	"github.com/Eugene7997/cheatcom/internal/tui"
)

var editCmd = &cobra.Command{
	Use:          "edit <id>",
	Short:        "Edit a cheat by ID",
	Args:         cobra.ExactArgs(1),
	RunE:         runEdit,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	s, err := store.Load()
	if err != nil {
		return err
	}

	id := args[0]
	c, ok := s.Get(id)
	if !ok {
		return fmt.Errorf("not found: %s", id)
	}

	res, err := tui.RunEdit(c)
	if err != nil {
		return fmt.Errorf("edit: %w", err)
	}
	if !res.Submitted {
		return nil
	}

	if strings.TrimSpace(res.Command) == "" {
		return fmt.Errorf("command must not be empty")
	}

	c.Command = res.Command
	c.Description = res.Description
	c.Tags = res.Tags
	c.UpdatedAt = time.Now().UTC()

	if err := s.Update(c); err != nil {
		return err
	}
	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("updated %s\n", id)
	return nil
}
