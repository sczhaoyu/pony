package util

import (
	"bytes"
	"encoding/binary"
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
	binary.Read(b_buf, binary.BigEndian, &x)
	return int(x)
}
func ReadData(conn net.Conn) ([]byte, error) {
	var l int = 4
	data := make([]byte, l)
	for l > 0 {
		i, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
		l = l - i
	}
	l = ByteSliceToInt(data)
	for l > 0 {
		data = make([]byte, l)
		i, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
		l = l - i
	}
	return data, nil
}
