package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Eugene7997/cheatcom/internal/action"
	"github.com/Eugene7997/cheatcom/internal/store"
	"github.com/Eugene7997/cheatcom/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:          "chc",
	Short:        "chc — a personal command cheatsheet",
	Long:         "Save and retrieve CLI commands you never want to retype.",
	RunE:         runRoot,
	SilenceUsage: true,
}

var (
	flagCopy  bool
	flagRun   bool
	flagPrint bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagCopy, "copy", false, "copy selected command to clipboard (default)")
	rootCmd.PersistentFlags().BoolVar(&flagRun, "run", false, "execute selected command in shell")
	rootCmd.PersistentFlags().BoolVar(&flagPrint, "print", false, "print selected command to stdout")
	rootCmd.MarkFlagsMutuallyExclusive("copy", "run", "print")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	s, err := store.Load()
	if err != nil {
		return err
	}

	defaultAct := parseActionString(s.DefaultAction())
	res, err := tui.Run(s.All(), defaultAct)
	if err != nil {
		return fmt.Errorf("tui: %w", err)
	}
	if res.Cancelled {
		return nil
	}

	act := resolveAction(res.Action)
	return action.Dispatch(res.Cheat, act)
}

func resolveAction(fromTUI action.Action) action.Action {
	if flagRun {
		return action.Run
	}
	if flagPrint {
		return action.Print
	}
	if flagCopy {
		return action.Copy
	}
	return fromTUI
}

func parseActionString(s string) action.Action {
	switch s {
	case "run":
		return action.Run
	case "print":
		return action.Print
	default:
		return action.Copy
	}
}
