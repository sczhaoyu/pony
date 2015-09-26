package common

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
)

type RequestHeader struct {
	SessionId string `json:"sessionId"`
	Command   string `json:"command"`
	Err       string `json:"err"`
}

type Request struct {
	Head *RequestHeader `json:"head"`
	Body []byte         `json:"body"`
}

func NewRequestJson(data []byte, sessionId string) []byte {
	var r Request
	r.Head = new(RequestHeader)
	r.Head.SessionId = sessionId
	r.Body = data
	return util.GetJsonByteLen(r)
}
func (r *Request) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
