package mdcoach

import (
	"errors"
	"fmt"
	"io/fs"
)

// ErrEmoteNotFound indicates that no configured filesystem contains an emote asset.
var ErrEmoteNotFound = fmt.Errorf("emote not found: %w", fs.ErrNotExist)

var _ fs.FS = (*emoteFilesystem)(nil)

type emoteFilesystem struct {
	filesystems []fs.FS
}

// NewEmoteFilesystem returns an fs.FS that searches each provided filesystem in order.
func NewEmoteFilesystem(filesystems ...fs.FS) fs.FS {
	return &emoteFilesystem{
		filesystems: append([]fs.FS(nil), filesystems...),
	}
}

func (e *emoteFilesystem) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	for _, filesystem := range e.filesystems {
		if filesystem == nil {
			continue
		}

		file, err := filesystem.Open(name)
		if err == nil {
			return file, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: ErrEmoteNotFound}
}
