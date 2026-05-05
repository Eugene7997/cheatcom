package action

import (
	"fmt"

	"github.com/Eugene7997/cheatcom/internal/clipboard"
	"github.com/Eugene7997/cheatcom/internal/runner"
	"github.com/Eugene7997/cheatcom/internal/store"
)

type Action int

const (
	Copy Action = iota
	Run
	Print
)

func Parse(s string) (Action, error) {
	switch s {
	case "copy":
		return Copy, nil
	case "run":
		return Run, nil
	case "print":
		return Print, nil
	default:
		return 0, fmt.Errorf("action must be one of: copy, run, print")
	}
}

func Dispatch(c store.Cheat, a Action) error {
	switch a {
	case Copy:
		if err := clipboard.Copy(c.Command); err != nil {
			return fmt.Errorf("copy to clipboard: %w", err)
		}
		return nil
	case Run:
		return runner.Run(c.Command)
	case Print:
		fmt.Println(c.Command)
		return nil
	default:
		return fmt.Errorf("unknown action %d", a)
	}
}
