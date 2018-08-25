package disk

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

type reader struct {
	img io.Reader
	buf bytes.Buffer
}

func NewReader(rgbaReader io.Reader) (io.Reader, error) {
	r := bufio.NewReader(&reader{img: rgbaReader})
	if _, err := r.ReadBytes(0x78); err != nil {
		return nil, errors.New("invalid disk image")
	}
	r.UnreadByte()

	solution, err := zlib.NewReader(r)
	if err != nil {
		return nil, errors.New("invalid disk image")
	}

	return solution, nil
}

func (r *reader) Read(b []byte) (n int, err error) {
	for r.buf.Len() < len(b) {
		var next [32]byte
		nn, err := r.img.Read(next[:])
		if err != nil {
			break
		}

		var data uint32
		var count uint
		for start := 3; start < nn; start += 4 {
			data >>= 3
			data |= uint32(next[start-3]&1) << 21
			data |= uint32(next[start-2]&1) << 22
			data |= uint32(next[start-1]&1) << 23
			count += 3
		}

		data >>= (24 - count)

		for i := count; i > 0; i -= 8 {
			r.buf.WriteByte(byte(data))
			data >>= 8
		}
	}

	return r.buf.Read(b)
}
