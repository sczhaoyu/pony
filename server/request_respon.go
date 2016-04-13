package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/util"
	"net"
	"time"
)

type RequestHeader struct {
	RequestId   string `json:"requestId"`          //唯一请求标示
	UserAddr    string `json:"userAddr"`           //用户客户端地址
	UserUid     string `json:"userUid"`            //客户端ID
	FaceCode    int    `json:"faceCode,omitempty"` //服务器需要处理的功能
	RequestTime int64  `json:"requestTime"`        //收到请求的时间

}

//服务器之间通讯的请求
type Request struct {
	Header RequestHeader `json:"header"` //头信息
	Body   interface{}   `json:"body"`   //客户端请求的数据
}

//反序
func (r *Request) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

//序列化json
func (r *Request) Marshal() []byte {
	data, _ := json.Marshal(r)
	return data
}
func NewRequest(conn net.Conn, data interface{}, s *Server) *Request {
	var r Request
	r.Header.UserAddr = conn.LocalAddr().String()
	r.Header.RequestId = util.GetUUID()
	r.Header.RequestTime = time.Now().Local().Unix()
	r.Body = data
	if s != nil {
		r.Header.UserUid = s.Id
	}
	return &r
}

type ResponseHeader struct {
	RequestId  string `json:"requestId"`          //唯一请求标示
	ResponsId  string `json:"responsId"`          //唯一响应标示
	UserAddr   string `json:"userAddr"`           //用户客户端地址
	UserUid    string `json:"userUid"`            //客户端ID
	FaceCode   int    `json:"faceCode,omitempty"` //服务器需要处理的功能
	ErrMsg     string `json:"errMsg,omitempty"`   //系统错误消息
	ErrCode    int    `json:"errCode"`            //错误代码
	ResponTime int64  `json:"responTime"`         //服务器处理完成时间
}

//服务器之间通讯的响应
type Respon struct {
	Header ResponseHeader `json:"header"` //头信息
	Body   interface{}    `json:"body"`   //响应数据
}

//反序
func (r *Respon) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

//序列化json
func (r *Respon) Marshal() []byte {
	data, _ := json.Marshal(r)
	return data
}
