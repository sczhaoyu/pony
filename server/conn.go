package server

import (
	"encoding/json"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/sczhaoyu/pony/errcode"
	"github.com/sczhaoyu/pony/util"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	Request *Request               //请求
	values  map[string]interface{} //请求的参数
	Body    []byte                 //请求的BODY
	Server  *Server                //服务器
	Session *SessionStore          //session会话
}

//获取请求的参数
func (c *Conn) Value(key string) *Value {
	var v Value
	v.val = c.values[key]
	return &v
}
func (c *Conn) Unmarshal(b interface{}) error {
	if len(c.Body) > 0 {
		err := json.Unmarshal(c.Body, b)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("body byte nil?")
}

//创建一个服务器端接收请求的链接
func NewConn(conn net.Conn, srv *Server, session *SessionStore, data []byte) (*Conn, error) {
	var c Conn
	c.Conn = conn
	c.Server = srv

	c.Session = session
	var r Request
	err := r.Unmarshal(data)
	if err != nil {
		return &c, err
	}
	c.Request = &r
	//取出data里面的body
	j, _ := simplejson.NewJson(data)
	body, err := j.Get("body").MarshalJSON()
	if err == nil {
		c.Body = body
	}
	//将BODY里面的内容转换
	m, err := j.Get("body").Map()
	if err == nil {
		c.values = m
	}
	//取出body后遍历里面的key val
	return &c, nil
}

//回应客户端json
func (c *Conn) WriteJson(b interface{}) {
	code := 0
	if c.Request != nil {
		code = c.Request.Header.FaceCode
	}
	c.WriteJsonFaceCode(b, code)
}
func (c *Conn) WriteJsonFaceCode(b interface{}, faceCode int) {
	var rsp Respon
	//请求数据的ID
	if c != nil && c.Request != nil {
		rsp.Header.RequestId = c.Request.Header.RequestId
		rsp.Header.UserUid = c.Request.Header.UserUid
		rsp.Header.UserAddr = c.Request.Header.UserAddr
	}
	rsp.Header.FaceCode = faceCode
	rsp.Header.ResponsId = util.GetUUID()
	rsp.Header.ResponTime = time.Now().Unix()
	switch e := b.(type) {
	case *errcode.Error:
		rsp.Header.ErrMsg = e.Msg
		rsp.Header.ErrCode = e.Code
	case error:
		rsp.Header.ErrMsg = e.Error()
		rsp.Header.ErrCode = -1
	default:
		rsp.Header.ErrCode = 0
		rsp.Body = b
	}
	d, _ := json.Marshal(&rsp)
	d = util.ByteLen(d)
	//加入包推送信息
	if c.Server != nil && c.Server.DPM != nil && rsp.Header.FaceCode != 520 && rsp.Header.FaceCode != 200 {
		c.Server.DPM.AddPkg(rsp.Header.ResponsId, c.Conn.RemoteAddr().String(), rsp.Header.UserUid, d)
	}

	c.Write(d)
}

type Value struct {
	val interface{}
}

func (v *Value) ToString() string {
	if v.val == nil {
		return ""
	}
	return util.ProphesyVal(v.val)
}
func (v *Value) ToInt() int {
	return util.InsToInt(v.val)
}
func (v *Value) ToInt64() int64 {
	return int64(util.InsToInt(v.val))
}
