package epub

import (
	"archive/zip"
	"io"
	"testing"
)

func readEPUBEntry(t *testing.T, file *zip.File) []byte {
	t.Helper()

	reader, err := file.Open()
	if err != nil {
		t.Fatalf("open EPUB entry %q: %v", file.Name, err)
	}
	content, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr != nil {
		t.Fatalf("read EPUB entry %q: %v", file.Name, readErr)
	}
	if closeErr != nil {
		t.Fatalf("close EPUB entry %q: %v", file.Name, closeErr)
	}
	return content
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}
