package socket

import (
	"errors"
	"io"
	"net"
)

func read(c net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 1024)

	for {
		tmp := make([]byte, 0, 1024)
		n, err := c.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		buf = append(buf, tmp[:n]...)
	}

	return buf, nil
}

func readString(c net.Conn) (string, error) {
	b, err := read(c)
	return string(b), err
}
