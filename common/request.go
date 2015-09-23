package common

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
)

type RequestHeader struct {
	SessionId string `json:"sessionId"`
	UserId    int64  `json:"userId"`
	FaceCode  int    `json:"faceCode"`
	Token     string `json:"token"`
}

type Request struct {
	Head *RequestHeader `json:"head"`
	Body []byte         `json:"body"`
}

func NewRequestJson(data []byte, sessionId string) []byte {
	var r Request
	r.Unmarshal(data)
	r.Head.SessionId = sessionId
	return util.GetJsonByteLen(r)
}
func (r *Request) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
