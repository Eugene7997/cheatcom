package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Eugene7997/cheatcom/internal/store"
)

var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "Manage chc configuration",
	SilenceUsage: true,
}

var configDefaultCmd = &cobra.Command{
	Use:          "default <copy|run|print>",
	Short:        "Set the default Enter action in the TUI",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         runConfigDefault,
}

func init() {
	configCmd.AddCommand(configDefaultCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigDefault(cmd *cobra.Command, args []string) error {
	a := args[0]
	if a != "copy" && a != "run" && a != "print" {
		return fmt.Errorf("action must be one of: copy, run, print")
	}
	s, err := store.Load()
	if err != nil {
		return err
	}
	if err := s.SetDefaultAction(a); err != nil {
		return err
	}
	fmt.Printf("Default action set to %q\n", a)
	return nil
}
