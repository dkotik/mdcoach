package picture

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func testdataServer(t *testing.T) string {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	// t.Cleanup(func() { // shutdown closes the connection
	// 	if err := listener.Close(); err != nil {
	// 		t.Fatal(err)
	// 	}
	// })

	s := &http.Server{
		Handler: http.FileServer(http.Dir("testdata")),
	}
	go s.Serve(listener)
	t.Cleanup(func() {
		if err := s.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
	return listener.Addr().String()
}

func TestPictureDownloadedFromInternet(t *testing.T) {
	if testing.Short() {
		t.Skip("accessing the Internet is slow")
	}

	p, err := NewInternetProvider(
		WithDestinationPath(t.TempDir()),
	)
	if err != nil {
		t.Fatal(err)
	}

	address := testdataServer(t)
	time.Sleep(time.Millisecond * 600)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	set, err := p.GetSourceSet(ctx, `http://`+address+`/notfound.jpg`)
	if err != nil {
		t.Fatal(err)
	}
	p.local.FinishScaling()
	spew.Dump(set)
}
