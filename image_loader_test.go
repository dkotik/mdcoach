package mdcoach

import "testing"

func TestLocalURLPath(t *testing.T) {
	url := newURLFromLocation(".")
	if url == nil {
		t.Fatal("expected non-nil URL")
	}
	if url.Host != "" {
		t.Fatal("expected local URL")
	}
	if url.Path != "." {
		t.Fatal("expected path to be '.'")
	}
}

func TestImageLoader(t *testing.T) {

}
