package codec

import (
	"io"
)

type I interface {
	MarshalBinary(w io.Writer)
	UnmarshalBinary(r io.Reader) (err error)
}
