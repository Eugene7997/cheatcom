package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempStore(t *testing.T) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cheats.yaml")
	s, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}
	return s, path
}

func TestLoadMissingFile(t *testing.T) {
	s, _ := tempStore(t)
	if got := len(s.All()); got != 0 {
		t.Fatalf("expected 0 cheats, got %d", got)
	}
}

func TestAddAndGet(t *testing.T) {
	s, _ := tempStore(t)
	c := Cheat{ID: "abc12345", Command: "ls -la", CreatedAt: time.Now().UTC()}
	if err := s.Add(c); err != nil {
		t.Fatal(err)
	}
	got, ok := s.Get("abc12345")
	if !ok {
		t.Fatal("expected to find cheat")
	}
	if got.Command != "ls -la" {
		t.Fatalf("unexpected command %q", got.Command)
	}
}

func TestDuplicateAdd(t *testing.T) {
	s, _ := tempStore(t)
	c := Cheat{ID: "dup12345", Command: "echo hi", CreatedAt: time.Now().UTC()}
	if err := s.Add(c); err != nil {
		t.Fatal(err)
	}
	if err := s.Add(c); err == nil {
		t.Fatal("expected error on duplicate ID")
	}
}

func TestSaveAndLoad(t *testing.T) {
	s, path := tempStore(t)
	c := Cheat{
		ID:          "sv123456",
		Command:     "git log --oneline",
		Description: "compact log",
		Tags:        []string{"git"},
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
	}
	if err := s.Add(c); err != nil {
		t.Fatal(err)
	}
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	s2, err := LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	cheats := s2.All()
	if len(cheats) != 1 {
		t.Fatalf("expected 1 cheat, got %d", len(cheats))
	}
	if cheats[0].Command != c.Command {
		t.Fatalf("command mismatch: %q vs %q", cheats[0].Command, c.Command)
	}
	if len(cheats[0].Tags) != 1 || cheats[0].Tags[0] != "git" {
		t.Fatalf("tags mismatch: %v", cheats[0].Tags)
	}
}

func TestDelete(t *testing.T) {
	s, _ := tempStore(t)
	c := Cheat{ID: "del12345", Command: "rm -rf /tmp/test", CreatedAt: time.Now().UTC()}
	_ = s.Add(c)

	if err := s.Delete("del12345"); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Get("del12345"); ok {
		t.Fatal("expected cheat to be gone")
	}

	if err := s.Delete("del12345"); err == nil {
		t.Fatal("expected error deleting non-existent id")
	}
}

func TestUpdate(t *testing.T) {
	s, _ := tempStore(t)
	c := Cheat{ID: "upd12345", Command: "old command", CreatedAt: time.Now().UTC()}
	_ = s.Add(c)

	c.Command = "new command"
	if err := s.Update(c); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Get("upd12345")
	if got.Command != "new command" {
		t.Fatalf("update did not stick: %q", got.Command)
	}
}

func TestAtomicSaveNoTempLeft(t *testing.T) {
	s, path := tempStore(t)
	_ = s.Add(Cheat{ID: "atm12345", Command: "echo atomic", CreatedAt: time.Now().UTC()})
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	dir := filepath.Dir(path)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if matched, _ := filepath.Match("cheats-*.yaml.tmp", e.Name()); matched {
			t.Fatalf("temp file left behind: %s", e.Name())
		}
	}
}

func TestMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cheats.yaml")
	if err := os.WriteFile(path, []byte("not: valid: yaml: [{{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFrom(path)
	if err == nil {
		t.Fatal("expected error on malformed YAML")
	}
}
