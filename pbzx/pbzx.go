// Package pbzx provides a reader for working with pbzx streams in Apple's
//.xip archive
// credit to Michael Lynn for the original work in python:
// https://gist.github.com/pudquick/ff412bcb29c9c1fa4b8d
package pbzx

import (
	"encoding/binary"
	"errors"
	"io"
)

// Copy copies from src to dst until either EOF is reached on src or an error
// occurs. It returns the number of bytes copied and the first error
// encountered while copying, if any.
// The Reader must be a pbzx encoded stream, such a xar Contents file from
// a .xip archive.
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	var intro introHeader
	if err := binary.Read(src, binary.BigEndian, &intro); err != nil {
		return 0, err // TODO wrap in custom error
	}
	if intro.Magic&0xffffffff != 0x70627a78 {
		return 0, errors.New("src not a pbzx stream")
	}
	for (intro.Flags & (1 << 24)) != 0 {
		var tag header
		if err := binary.Read(src, binary.BigEndian, &tag); err != nil {
			return 0, err
		}
		intro.Flags = tag.Flags
		n, err := io.CopyN(dst, src, int64(tag.Size))
		if err != nil {
			return written + n, err
		}
		written += n
	}
	return written, nil
}

// introHeader is the first header, indicating a pbzx stream.
type introHeader struct {
	Magic uint32
	Flags uint64
}

// header is a standard header of an xz encoded chunk.
type header struct {
	Flags uint64
	Size  uint64
}
