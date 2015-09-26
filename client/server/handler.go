package server

import (
	"github.com/sczhaoyu/pony/common"
)

func handler(lc *LogicConn, data []byte) {
	var rsp common.Response
	err := rsp.Unmarshal(data)
	if err == nil {
		switch rsp.Head.Command {
		case common.RADIO:
			s := lc.LSM.ClientServer.Session.Session
			for k, _ := range s {
				r := common.NewResponseSid(k, rsp.Body)
				lc.LSM.RspChan <- r
			}
		default:
			lc.LSM.RspChan <- &rsp
		}

	}

}
