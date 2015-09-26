package common

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
)

type ResponseHead struct {
	Err       string `json:"err"`
	SessionId string `json:"sessionId"`
	Command   string `json:"command"`
}
type Response struct {
	Head *ResponseHead `json:"head"`
	Body interface{}   `json:"body"`
}

func AuthResponse(command string, data interface{}) *Response {
	var r Response
	r.Head = new(ResponseHead)
	r.Head.Command = command
	r.Body = data
	return &r
}
func NewResponseSid(sessionId string, data interface{}) *Response {
	var r Response
	r.Head = new(ResponseHead)
	r.Head.SessionId = sessionId
	r.Body = data
	return &r
}
func (r *Response) GetJsonByte() []byte {
	data, _ := json.Marshal(r)
	return data
}
func (r *Response) GetJson() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return util.ByteLen(data)
}
func (r *Response) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
