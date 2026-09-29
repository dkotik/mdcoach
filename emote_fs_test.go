package mdcoach

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestEmoteFilesystem(t *testing.T) {
	tests := []struct {
		name        string
		filesystems []fs.FS
		path        string
		wantContent string
		wantErrors  []error
	}{
		{
			name: "uses first filesystem containing the file",
			filesystems: []fs.FS{
				fstest.MapFS{"assets/smile.png": &fstest.MapFile{Data: []byte("first")}},
				fstest.MapFS{"assets/smile.png": &fstest.MapFile{Data: []byte("second")}},
			},
			path:        "assets/smile.png",
			wantContent: "first",
		},
		{
			name: "falls back to later filesystems",
			filesystems: []fs.FS{
				fstest.MapFS{},
				fstest.MapFS{"assets/smile.png": &fstest.MapFile{Data: []byte("found")}},
			},
			path:        "assets/smile.png",
			wantContent: "found",
		},
		{
			name:        "reports emote not found after searching all filesystems",
			filesystems: []fs.FS{fstest.MapFS{}, fstest.MapFS{}},
			path:        "assets/missing.png",
			wantErrors:  []error{ErrEmoteNotFound, fs.ErrNotExist},
		},
		{
			name:       "rejects invalid paths",
			path:       "../smile.png",
			wantErrors: []error{fs.ErrInvalid},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := fs.ReadFile(NewEmoteFilesystem(tt.filesystems...), tt.path)
			if len(tt.wantErrors) > 0 {
				if err == nil {
					t.Fatal("expected an error")
				}
				for _, wantErr := range tt.wantErrors {
					if !errors.Is(err, wantErr) {
						t.Errorf("error = %v, want it to match %v", err, wantErr)
					}
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.wantContent {
				t.Fatalf("file contents = %q, want %q", data, tt.wantContent)
			}
		})
	}
}
