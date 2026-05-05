package store

import (
	"strings"
	"time"
)

type Cheat struct {
	ID          string    `yaml:"id"                    mapstructure:"id"`
	Command     string    `yaml:"command"               mapstructure:"command"`
	Description string    `yaml:"description,omitempty" mapstructure:"description"`
	Tags        []string  `yaml:"tags,omitempty"        mapstructure:"tags"`
	CreatedAt   time.Time `yaml:"created_at"            mapstructure:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at,omitempty"  mapstructure:"updated_at"`
}

func NormalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
