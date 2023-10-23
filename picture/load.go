package picture

import (
	"image"
	"os"

	_ "image/gif" // add GIF format to image reader
	_ "image/png" // add PNG format to image reader
)

// LoadImage reads contents of a file into [image.Image].
func LoadImage(p string) (m image.Image, err error) {
	handle, err := os.Open(p)
	if err != nil {
		return m, err
	}
	defer handle.Close()
	m, _, err = image.Decode(handle)
	return m, err
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
