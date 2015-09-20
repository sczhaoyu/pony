package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func IntToByteSlice(i int) []byte {
	var x int32 = int32(i)
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, x)
	return b_buf.Bytes()
}
func ByteSliceToInt(data []byte) int {
	b_buf := bytes.NewBuffer(data)
	var x int32
	err := binary.Read(b_buf, binary.BigEndian, &x)
	if err != nil {
		return 0
	}
	return int(x)
}
func ReadData(conn net.Conn, length int) ([]byte, error) {
	var l int = 4
	data := make([]byte, l)
	var i int = 0
	var err error
	for l > 0 {
		i, err = conn.Read(data)
		if err != nil {
			return nil, err
		}
		l = l - i
	}
	l = ByteSliceToInt(data)
	if l == 0 {
		return nil, errors.New("data nil")
	}
	if l > length {
		return nil, errors.New("data must  < " + fmt.Sprintf("%d", length))
	}
	data = make([]byte, l)
	for l > 0 {
		i, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
		l = l - i
	}
	return data, nil
}
