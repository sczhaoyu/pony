package server

import (
	simplejson "github.com/bitly/go-simplejson"
	"github.com/sczhaoyu/pony/common"
	"log"
)

func handler(lc *LogicConn, data []byte) {
	sj, err := simplejson.NewJson(data)
	if err != nil {
		return
	}
	var head common.ResponseHead
	var body []byte
	h, hr := sj.Get("head").MarshalJSON()
	if hr == nil {
		head.Unmarshal(h)
	}
	b, br := sj.Get("body").MarshalJSON()
	if br == nil {
		body = b
	}
	log.Println(string(b))
	switch head.Command {
	case common.RADIO:
		s := lc.LSM.ClientServer.Session.Session
		for k, _ := range s {
			var r Rsp
			r.SessionId = k
			r.Data = body
			lc.LSM.RspChan <- &r
		}
	default:
		var r Rsp
		r.SessionId = head.SessionId
		r.Data = body
		lc.LSM.RspChan <- &r
	}

}
