package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Eugene7997/cheatcom/internal/store"
)

// storeDir redirects APPDATA to a fresh temp dir so store.Load() uses an
// isolated cheats.yaml. Returns the path to that file.
func storeDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	// Suppress cobra's own error printing so test output stays clean.
	rootCmd.SetErr(io.Discard)
	t.Cleanup(func() { rootCmd.SetErr(nil) })
	return filepath.Join(dir, "cheatcom", "cheats.yaml")
}

// seedCheat writes a cheat directly into the store file so commands have
// pre-existing data to operate on.
func seedCheat(t *testing.T, path string, c store.Cheat) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	s, err := store.LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Add(c); err != nil {
		t.Fatal(err)
	}
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
}

// exec sets cobra args and runs the root command, returning any error.
func exec(args ...string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

// captureStdout replaces os.Stdout with a pipe, runs fn, restores os.Stdout,
// and returns what was written. Needed for commands that write directly to
// os.Stdout rather than cmd.OutOrStdout().
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fnErr := fn()

	w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String(), fnErr
}

func TestAddCmd(t *testing.T) {
	path := storeDir(t)
	addDescription = ""
	addTags = nil

	if err := exec("add", "git log --oneline", "-d", "compact log", "-t", "git"); err != nil {
		t.Fatal(err)
	}

	s, err := store.LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	cheats := s.All()
	if len(cheats) != 1 {
		t.Fatalf("expected 1 cheat, got %d", len(cheats))
	}
	c := cheats[0]
	if c.Command != "git log --oneline" {
		t.Fatalf("command: got %q", c.Command)
	}
	if c.Description != "compact log" {
		t.Fatalf("description: got %q", c.Description)
	}
	if len(c.Tags) != 1 || c.Tags[0] != "git" {
		t.Fatalf("tags: got %v", c.Tags)
	}
	if c.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestAddCmdNoFlags(t *testing.T) {
	path := storeDir(t)
	addDescription = ""
	addTags = nil

	if err := exec("add", "echo hello"); err != nil {
		t.Fatal(err)
	}

	s, err := store.LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	cheats := s.All()
	if len(cheats) != 1 {
		t.Fatalf("expected 1 cheat, got %d", len(cheats))
	}
	if cheats[0].Command != "echo hello" {
		t.Fatalf("command: got %q", cheats[0].Command)
	}
	if cheats[0].Description != "" {
		t.Fatalf("expected empty description, got %q", cheats[0].Description)
	}
	if len(cheats[0].Tags) != 0 {
		t.Fatalf("expected no tags, got %v", cheats[0].Tags)
	}
}

func TestAddCmdMissingArg(t *testing.T) {
	storeDir(t)
	addDescription = ""
	addTags = nil

	if err := exec("add"); err == nil {
		t.Fatal("expected error when command argument is missing")
	}
}

func TestAddCmdMultipleCheats(t *testing.T) {
	path := storeDir(t)
	addDescription = ""
	addTags = nil

	if err := exec("add", "ls -la"); err != nil {
		t.Fatal(err)
	}
	addDescription = ""
	addTags = nil
	if err := exec("add", "pwd"); err != nil {
		t.Fatal(err)
	}

	s, err := store.LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.All()) != 2 {
		t.Fatalf("expected 2 cheats, got %d", len(s.All()))
	}
}

func TestListCmdEmpty(t *testing.T) {
	storeDir(t)
	out, err := captureStdout(t, func() error { return exec("list") })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No cheats yet") {
		t.Fatalf("expected empty-store message, got %q", out)
	}
}

func TestListCmdWithCheats(t *testing.T) {
	path := storeDir(t)
	seedCheat(t, path, store.Cheat{
		ID:          "listid01",
		Command:     "ls -la",
		Description: "list files",
		Tags:        []string{"shell"},
		CreatedAt:   time.Now().UTC(),
	})

	out, err := captureStdout(t, func() error { return exec("list") })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "listid01") {
		t.Fatalf("expected ID in output, got %q", out)
	}
	if !strings.Contains(out, "ls -la") {
		t.Fatalf("expected command in output, got %q", out)
	}
	if !strings.Contains(out, "#shell") {
		t.Fatalf("expected tag in output, got %q", out)
	}
}

func TestRmCmd(t *testing.T) {
	path := storeDir(t)
	seedCheat(t, path, store.Cheat{
		ID:        "rmid0001",
		Command:   "docker ps",
		CreatedAt: time.Now().UTC(),
	})

	if err := exec("rm", "rmid0001"); err != nil {
		t.Fatal(err)
	}

	s, err := store.LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.All()) != 0 {
		t.Fatal("expected empty store after rm")
	}
}

func TestRmCmdNotFound(t *testing.T) {
	storeDir(t)
	if err := exec("rm", "notexist"); err == nil {
		t.Fatal("expected error when ID not found")
	}
}

func TestConfigDefaultCmd(t *testing.T) {
	path := storeDir(t)

	for _, action := range []string{"copy", "run", "print"} {
		t.Run(action, func(t *testing.T) {
			if err := exec("config", "default", action); err != nil {
				t.Fatalf("config default %s: %v", action, err)
			}
			s, err := store.LoadFrom(path)
			if err != nil {
				t.Fatal(err)
			}
			if got := s.DefaultAction(); got != action {
				t.Fatalf("DefaultAction: got %q, want %q", got, action)
			}
		})
	}
}

func TestConfigDefaultCmdInvalid(t *testing.T) {
	storeDir(t)
	if err := exec("config", "default", "invalid"); err == nil {
		t.Fatal("expected error for invalid action")
	}
}
