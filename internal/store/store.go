package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Store struct {
	v             *viper.Viper
	path          string
	cheats        []Cheat
	defaultAction string
}

type fileFormat struct {
	Version       int     `yaml:"version"`
	DefaultAction string  `yaml:"default_action,omitempty"`
	Cheats        []Cheat `yaml:"cheats"`
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config dir: %w", err)
	}
	return filepath.Join(dir, "cheatcom", "cheats.yaml"), nil
}

func Load() (*Store, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

func LoadFrom(path string) (*Store, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	s := &Store{v: v, path: path}

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return s, nil
		}
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return s, nil
		}
		return nil, fmt.Errorf("cheats file is corrupted at %s — fix manually or move it aside: %w", path, err)
	}

	if err := v.UnmarshalKey("cheats", &s.cheats); err != nil {
		return nil, fmt.Errorf("cheats file is corrupted at %s — fix manually or move it aside: %w", path, err)
	}
	s.defaultAction = v.GetString("default_action")
	return s, nil
}

func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := yaml.Marshal(fileFormat{Version: 1, DefaultAction: s.defaultAction, Cheats: s.cheats})
	if err != nil {
		return fmt.Errorf("marshal cheats: %w", err)
	}

	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, "cheats-*.yaml.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()

	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("sync temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err = os.Rename(tmpName, s.path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("rename to final path: %w", err)
	}
	return nil
}

func (s *Store) All() []Cheat {
	out := make([]Cheat, len(s.cheats))
	copy(out, s.cheats)
	return out
}

func (s *Store) Get(id string) (Cheat, bool) {
	for _, c := range s.cheats {
		if c.ID == id {
			return c, true
		}
	}
	return Cheat{}, false
}

func (s *Store) Add(c Cheat) error {
	if _, exists := s.Get(c.ID); exists {
		return fmt.Errorf("id %q already exists", c.ID)
	}
	s.cheats = append(s.cheats, c)
	return nil
}

func (s *Store) Delete(id string) error {
	for i, c := range s.cheats {
		if c.ID == id {
			s.cheats = append(s.cheats[:i], s.cheats[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cheat %q not found", id)
}

func (s *Store) Update(updated Cheat) error {
	for i, c := range s.cheats {
		if c.ID == updated.ID {
			s.cheats[i] = updated
			return nil
		}
	}
	return fmt.Errorf("cheat %q not found", updated.ID)
}

// DefaultAction returns the configured default action string ("copy", "run", "print").
// Returns "copy" if none is set.
func (s *Store) DefaultAction() string {
	if s.defaultAction == "" {
		return "copy"
	}
	return s.defaultAction
}

func (s *Store) NewUniqueID() (string, error) {
	for range 10 {
		id := NewID()
		if _, exists := s.Get(id); !exists {
			return id, nil
		}
	}
	return "", fmt.Errorf("could not generate a unique ID after 10 attempts")
}

func (s *Store) SetDefaultAction(a string) error {
	s.defaultAction = a
	return s.Save()
}
