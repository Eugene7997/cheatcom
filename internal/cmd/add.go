package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Eugene7997/cheatcom/internal/store"
)

var addCmd = &cobra.Command{
	Use:          "add <command>",
	Short:        "Add a command to the cheatsheet",
	Args:         cobra.ExactArgs(1),
	RunE:         runAdd,
	SilenceUsage: true,
}

var (
	addDescription string
	addTags        []string
)

func init() {
	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Short description")
	addCmd.Flags().StringSliceVarP(&addTags, "tags", "t", nil, "Comma-separated tags")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	s, err := store.Load()
	if err != nil {
		return err
	}

	var id string
	for range 3 {
		candidate := store.NewID()
		if _, exists := s.Get(candidate); !exists {
			id = candidate
			break
		}
	}
	if id == "" {
		return fmt.Errorf("could not generate a unique ID, please try again")
	}

	c := store.Cheat{
		ID:          id,
		Command:     args[0],
		Description: addDescription,
		Tags:        normalizeTags(addTags),
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.Add(c); err != nil {
		return err
	}
	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("added %s\n", id)
	return nil
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
