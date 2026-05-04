package tui

import (
	"strings"

	"github.com/Eugene7997/cheatcom/internal/store"
)

type cheatItem struct {
	c store.Cheat
}

func (i cheatItem) Title() string {
	return strings.ReplaceAll(i.c.Command, "\n", "↵")
}

func (i cheatItem) Description() string {
	parts := []string{}
	if i.c.Description != "" {
		parts = append(parts, i.c.Description)
	}
	for _, t := range i.c.Tags {
		parts = append(parts, "#"+t)
	}
	return strings.Join(parts, " · ")
}

func (i cheatItem) FilterValue() string {
	return i.c.Command + " " + i.c.Description + " " + strings.Join(i.c.Tags, " ")
}
