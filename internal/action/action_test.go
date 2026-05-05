package action

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Eugene7997/cheatcom/internal/store"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input   string
		want    Action
		wantErr bool
	}{
		{"copy", Copy, false},
		{"run", Run, false},
		{"print", Print, false},
		{"", 0, true},
		{"invalid", 0, true},
		{"Copy", 0, true},
		{"RUN", 0, true},
	}
	for _, tc := range cases {
		name := tc.input
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			got, err := Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q): unexpected error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("Parse(%q): got %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestDispatchPrint(t *testing.T) {
	// Capture os.Stdout because Dispatch(Print) calls fmt.Println directly.
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	c := store.Cheat{Command: "echo hello world"}
	dispatchErr := Dispatch(c, Print)

	w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	if dispatchErr != nil {
		t.Fatalf("unexpected error: %v", dispatchErr)
	}
	if !strings.Contains(buf.String(), "echo hello world") {
		t.Fatalf("expected output to contain %q, got %q", "echo hello world", buf.String())
	}
}

func TestDispatchUnknownAction(t *testing.T) {
	c := store.Cheat{Command: "echo"}
	err := Dispatch(c, Action(99))
	if err == nil {
		t.Fatal("expected error for unknown action, got nil")
	}
}
