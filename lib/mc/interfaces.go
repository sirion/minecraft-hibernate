package mc

import "io"

type Reader interface {
	io.Reader
	io.ByteReader
}
