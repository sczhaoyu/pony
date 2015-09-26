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

func (r *Response) GetJson() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return util.ByteLen(data)
}
