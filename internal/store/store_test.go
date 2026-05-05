package store

import (
	"os"
	"path/filepath"
	"strings"
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

func TestUpdateNotFound(t *testing.T) {
	s, _ := tempStore(t)
	err := s.Update(Cheat{ID: "missing1", Command: "echo"})
	if err == nil {
		t.Fatal("expected error when updating non-existent cheat")
	}
}

func TestNormalizeTags(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		want  []string
	}{
		{"nil input", nil, []string{}},
		{"empty slice", []string{}, []string{}},
		{"already lowercase", []string{"git"}, []string{"git"}},
		{"uppercase converted", []string{"Git", "DOCKER"}, []string{"git", "docker"}},
		{"whitespace trimmed", []string{"  git  ", " docker"}, []string{"git", "docker"}},
		{"empty strings dropped", []string{"git", "", "  "}, []string{"git"}},
		{"mixed", []string{"  Shell ", "GIT", "", "docker"}, []string{"shell", "git", "docker"}},
		{"duplicates kept", []string{"git", "git"}, []string{"git", "git"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeTags(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("length mismatch: got %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("index %d: got %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestNewID(t *testing.T) {
	id := NewID()
	if len(id) != 8 {
		t.Fatalf("expected 8-char ID, got %q (len %d)", id, len(id))
	}
	for _, ch := range id {
		if !strings.ContainsRune(idAlphabet, ch) {
			t.Fatalf("ID %q contains character %q outside alphabet", id, ch)
		}
	}
}

func TestNewIDUniqueness(t *testing.T) {
	seen := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		id := NewID()
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate ID generated: %q", id)
		}
		seen[id] = struct{}{}
	}
}

func TestDefaultAction(t *testing.T) {
	s, _ := tempStore(t)
	if got := s.DefaultAction(); got != "copy" {
		t.Fatalf("expected default %q, got %q", "copy", got)
	}
}

func TestSetDefaultAction(t *testing.T) {
	s, path := tempStore(t)

	for _, a := range []string{"run", "print", "copy"} {
		if err := s.SetDefaultAction(a); err != nil {
			t.Fatalf("SetDefaultAction(%q): %v", a, err)
		}
		if got := s.DefaultAction(); got != a {
			t.Fatalf("DefaultAction: got %q, want %q", got, a)
		}

		s2, err := LoadFrom(path)
		if err != nil {
			t.Fatalf("reload after SetDefaultAction(%q): %v", a, err)
		}
		if got := s2.DefaultAction(); got != a {
			t.Fatalf("reloaded DefaultAction: got %q, want %q", got, a)
		}
	}
}

func TestNewUniqueID(t *testing.T) {
	s, _ := tempStore(t)
	id, err := s.NewUniqueID()
	if err != nil {
		t.Fatalf("NewUniqueID: %v", err)
	}
	if len(id) != 8 {
		t.Fatalf("expected 8-char ID, got %q", id)
	}
	if _, exists := s.Get(id); exists {
		t.Fatal("generated ID should not already be in store")
	}
}
