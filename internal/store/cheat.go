package store

import "time"

type Cheat struct {
	ID          string    `yaml:"id"                    mapstructure:"id"`
	Command     string    `yaml:"command"               mapstructure:"command"`
	Description string    `yaml:"description,omitempty" mapstructure:"description"`
	Tags        []string  `yaml:"tags,omitempty"        mapstructure:"tags"`
	CreatedAt   time.Time `yaml:"created_at"            mapstructure:"created_at"`
}
