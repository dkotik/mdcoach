package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandUserPath(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "home directory",
			path: "~",
			want: home,
		},
		{
			name: "path under home directory",
			path: "~/Documents/notes.md",
			want: filepath.Join(home, "Documents", "notes.md"),
		},
		{
			name: "ordinary relative path",
			path: "notes.md",
			want: "notes.md",
		},
		{
			name: "tilde in later path component",
			path: "docs/~/notes.md",
			want: "docs/~/notes.md",
		},
		{
			name: "other user home notation",
			path: "~ada/notes.md",
			want: "~ada/notes.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := expandUserPath(tt.path)
			if err != nil {
				t.Fatalf("expandUserPath(%q): %v", tt.path, err)
			}
			if got != tt.want {
				t.Errorf("expandUserPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
