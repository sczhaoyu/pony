package util

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

//生成一个UUID
func GetUUID() string {
	return strings.Replace(uuid.NewUUID().String(), "-", "", -1)
}
func GetJsonByteLen(d interface{}) []byte {
	data, _ := json.Marshal(d)
	return ByteLen(data)
}
func ByteLen(data []byte) []byte {
	tmp := IntToByteSlice(len(data))
	return append(tmp, data...)
}
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
	if l > length || l <= 0 {
		return nil, errors.New("data must  < " + fmt.Sprintf("%d", length) + "> 0 !")
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
