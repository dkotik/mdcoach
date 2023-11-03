package picture

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/OneOfOne/xxhash"
)

type Loader interface {
	LoadImage(string) (hash string, m image.Image, err error)
}

// LoadImage reads contents of a file into [image.Image].
func LoadImage(p string) (hash string, m image.Image, err error) {
	r, err := os.Open(p)
	if err != nil {
		return "", m, err
	}
	defer func() {
		err = errors.Join(err, r.Close())
	}()

	h := xxhash.New64()
	if _, err = io.CopyN(h, r, 4096); err != nil {
		return "", m, err
	}
	if _, err = r.Seek(0, 0); err != nil {
		return "", m, err
	}

	m, _, err = image.Decode(r)
	return fmt.Sprintf("%x", h.Sum(nil)), m, err
}

func copyFile(ctx context.Context, destination, source string) error {
	err := isContextAlive(ctx)
	if err != nil {
		return err
	}
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, r.Close())
	}()

	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, w.Close())
	}()

	_, err = io.Copy(w, r)
	return err
}

// DownloadImage reads contents of a URL into [image.Image].
// func DownloadImage(uri string, local bool) (m image.Image, err error) {
// 	var handle io.ReadCloser
// 	if local {
// 		handle, err = os.Open(uri)
// 		if err != nil {
// 			return m, err
// 		}
// 	} else if httpError != nil {
// 		return nil, fmt.Errorf(`network access blocked by previous failure`)
// 	} else {
// 		resp, err2 := (&http.Client{Timeout: httpTimeout}).Get(uri)
// 		if err2 != nil {
// 			httpError = err2
// 			return m, httpError
// 		}
// 		handle = resp.Body
// 	}
// 	defer handle.Close()
// 	m, _, err = image.Decode(handle)
// 	return m, err
// }
