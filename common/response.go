package common

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
)

type ResponseHead struct {
	UserId    int64  `json:"userId"`
	Err       string `json:"err"`
	State     int    `json:"state"`
	SessionId string `json:"sessionId"`
	FaceCode  int    `json:"faceCode"`
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
