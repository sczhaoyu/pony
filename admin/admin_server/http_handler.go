package admin_server

import (
	"errors"
	"github.com/sczhaoyu/pony/common"
	"net/http"
	"sort"
)

type LSAddrList []*common.LSAddr

func (list LSAddrList) Len() int {
	return len(list)
}

func (list LSAddrList) Less(i, j int) bool {
	return list[i].Num < list[j].Num
}

func (list LSAddrList) Swap(i, j int) {
	var temp *common.LSAddr = list[i]
	list[i] = list[j]
	list[j] = temp
}

//获取逻辑服务器IP地址 返回链接最小的IP
func getLogicAddr(w http.ResponseWriter, r *http.Request) {
	var ret []*common.LSAddr
	//取出逻辑服务器IP比较
	for _, v := range admin.CS {
		if v.ServerType == common.LS {
			var l common.LSAddr
			l.Addr = v.Addr
			l.Num = v.ClientNum
			ret = append(ret, &l)
		}
	}
	if len(ret) == 0 {
		d := common.NewResponse(errors.New("not found logic server")).GetJsonByte()
		w.Write(d)
		return
	}
	//排序
	list := LSAddrList(ret)
	sort.Sort(list)
	w.Write(common.NewResponse(list[0:1]).GetJsonByte())
}
func getLogicLiost(w http.ResponseWriter, r *http.Request) {
	var ret []*common.LSAddr = make([]*common.LSAddr, 0, len(admin.CS))
	for _, v := range admin.CS {
		if v.ServerType == common.LS {
			var l common.LSAddr
			l.Addr = v.Addr
			l.Num = v.ClientNum
			ret = append(ret, &l)
		}
	}
	w.Write(common.NewResponse(ret).GetJsonByte())
}
