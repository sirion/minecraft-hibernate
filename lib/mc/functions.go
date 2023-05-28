package mc

import (
	"errors"
	"io"
)

func ReadData(cr io.Reader, length int) ([]byte, error) {
	buf := make([]byte, length)
	read, err := io.ReadFull(cr, buf)
	if err != nil {
		return nil, err
	}
	if read != length {
		return nil, errors.New("WTF? read != length")
	}
	return buf, nil
}
