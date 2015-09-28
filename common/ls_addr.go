package common

import (
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
)

type LSAddr struct {
	Addr string `json:"addr"` //链接的地址
	Num  int    `json:"num"`  //链接数量
}

func (l *LSAddr) Marshal() []byte {
	data, _ := json.Marshal(l)
	return data
}
func (l *LSAddr) Unmarshal(data []byte) {
	json.Unmarshal(data, l)
}
func UnmarshalLSAddr(response []byte) []LSAddr {
	sj, jerr := simplejson.NewJson(response)
	if jerr != nil {
		return nil
	}
	data, err := sj.Get("body").MarshalJSON()
	//创建链接
	var ret []LSAddr
	if err != nil {
		return nil
	}
	json.Unmarshal(data, &ret)
	return ret
}
func GetLSAddr(response []byte) *LSAddr {
	ret := UnmarshalLSAddr(response)
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}
func AutoLSAddrReq(command, addr string, num int) *Request {
	var la LSAddr
	la.Addr = addr
	la.Num = num
	req := AuthRequest(command, la.Marshal())
	return req
}
