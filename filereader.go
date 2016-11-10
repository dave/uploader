package uploader

import (
	"io"

	"github.com/gopherjs/gopherjs/js"
)

type FileReader struct {
	file   *js.Object
	reader *js.Object
	offset int
	size   int
}

func NewFileReader(file *js.Object) *FileReader {
	return &FileReader{file, js.Global.Get("FileReader").New(), 0, file.Get("size").Int()}
}

func (fr *FileReader) Read(p []byte) (int, error) {
	if fr.offset == fr.size {
		return 0, io.EOF
	}
	type result struct {
		sz  int
		err error
	}
	c := make(chan result)
	go func() {
		fr.reader.Set("onloadend", func(evt *js.Object) {
			arr := js.Global.Get("Uint8Array").New(fr.reader.Get("result"))
			buf := arr.Interface().([]byte)
			go func() {
				if len(buf) == 0 {
					c <- result{0, io.EOF}
				} else {
					copy(p, buf)
					c <- result{len(buf), nil}
				}
			}()
		})
	}()
	e := fr.offset + len(p)
	if e > fr.size {
		e = fr.size
	}
	blob := fr.file.Call("slice", fr.offset, e)
	fr.reader.Call("readAsArrayBuffer", blob)
	res := <-c
	fr.offset = e
	return res.sz, res.err
}
