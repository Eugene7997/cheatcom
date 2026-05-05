package cmd

import (
	"fmt"
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

	id, err := s.NewUniqueID()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	c := store.Cheat{
		ID:          id,
		Command:     args[0],
		Description: addDescription,
		Tags:        store.NormalizeTags(addTags),
		CreatedAt:   now,
		UpdatedAt:   now,
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
