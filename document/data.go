package document

// Base64 data object
// <img src="data:image/gif;base64,xxxxbase64encodedtextxxxx">

import "io"

type DataWriter struct {
	w io.Writer
}

func (d *DataWriter) Open(ID string) (err error) {
	_, err = io.WriteString(d.w, `<div style="display: none;"><div id="`)
	if err != nil {
		return err
	}
	_, err = io.WriteString(d.w, ID)
	if err != nil {
		return err
	}
	_, err = io.WriteString(d.w, `">`)
	return err
}

func (d *DataWriter) Add(data []byte) error {
	return writeEscapedHTML(d.w, string(data))
}

func (d *DataWriter) Close() (err error) {
	_, err = io.WriteString(d.w, `</div></div>`)
	return err
}

func NewDataWriter(w io.Writer, ID string) (*DataWriter, func() error, error) {
	dw := &DataWriter{w: w}
	err := dw.Open(ID)
	if err != nil {
		return nil, nil, err
	}
	return dw, dw.Close, nil
}
