package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/util"
)

type RequestHeader struct {
	UserId   int64  `json:"userId"`
	UserAddr string `json:"userAddr"`
	FaceCode int    `json:"faceCode"`
	Token    string `json:"token"`
	Cid      string `json:"cid"`
}

type Request struct {
	Head *RequestHeader `json:"head"`
	Body []byte         `json:"body"`
}

func NewRequest(conn *common.Conn, data []byte) *Request {
	var r Request
	json.Unmarshal(data, &r)
	r.Head.UserAddr = conn.RemoteAddr().String()
	r.Head.Cid = conn.UUID
	return &r
}
func (r *Request) GetJson() []byte {
	data, _ := json.Marshal(r)
	return util.ByteLen(data)
}
