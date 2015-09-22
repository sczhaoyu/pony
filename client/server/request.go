package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
)

type RequestHeader struct {
	UserId   int64  `json:"userId"`
	FaceCode int    `json:"faceCode"`
	Token    string `json:"token"`
}

type Request struct {
	Head *RequestHeader `json:"head"`
	Body []byte         `json:"body"`
}

func (r *Request) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *Request) GetJson() []byte {
	data, _ := json.Marshal(r)
	return util.ByteLen(data)
}
