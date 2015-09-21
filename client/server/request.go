package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
	"net"
)

type RequestHeader struct {
	UserId   int64  `json:"userId"`
	UserAddr string `json:"userAddr"`
	FaceCode int    `json:"faceCode"`
	Token    string `json:"token"`
}

type Request struct {
	Head *RequestHeader `json:"head"`
	Body []byte         `json:"body"`
}

func NewRequest(conn net.Conn, data []byte) *Request {
	var r Request
	json.Unmarshal(data, &r)
	r.Head.UserAddr = conn.RemoteAddr().String()
	return &r
}
func (r *Request) GetJson() []byte {
	data, _ := json.Marshal(r)
	return util.ByteLen(data)
}
