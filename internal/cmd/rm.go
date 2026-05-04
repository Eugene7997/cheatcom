package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Eugene7997/cheatcom/internal/store"
)

var rmCmd = &cobra.Command{
	Use:          "rm <id>",
	Short:        "Remove a cheat by ID",
	Args:         cobra.ExactArgs(1),
	RunE:         runRm,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func runRm(cmd *cobra.Command, args []string) error {
	s, err := store.Load()
	if err != nil {
		return err
	}

	id := args[0]
	if err := s.Delete(id); err != nil {
		return fmt.Errorf("not found: %s", id)
	}

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("removed %s\n", id)
	return nil
}
