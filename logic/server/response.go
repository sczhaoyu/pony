package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
)

type ResponseHead struct {
	UserId int64  `json:"userId"`
	Uuid   string `json:"uuid"`
	Addr   string `json:"addr"`
	Msg    string `json:"msg"`
	State  int    `json:"state"`
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
