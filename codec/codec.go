package codec

import (
	"io"
)

type I interface {
	MarshalWrite(w io.Writer)
	UnmarshalRead(r io.Reader) (err error)
}
