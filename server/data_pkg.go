package server

import (
	"log"
	"sync"
	"time"
)

type DataPkg struct {
	MessageId string //消息包的唯一标示
	State     bool   //false 未确认，true已经确认
	Addr      string //发送目标地址
	Data      []byte //数据包
	At        int64  //发送时间
	Uid       string //用户ID
}

//数据收发确认补包管理
type DataPkgManager struct {
	Mx          sync.RWMutex        //操作锁
	Pkgs        map[string]*DataPkg //消息包
	GCTime      time.Duration       //GC的时间
	SendTime    time.Duration       //数据包重新发送时间
	ReceiveMsgs map[string]bool     //确定收到的包
	CH          chan *DataPkg       //数据写出通道(需要重新发包)
	CloseCH     chan int            //关闭
}

//重置发送数据包
func (d *DataPkgManager) ResetSend() {
	t := time.NewTicker(d.SendTime * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			d.Mx.Lock()
			for k, v := range d.Pkgs {
				//没有确认的消息包
				if v.State == false {
					//检测时间，是否重新推送
					if (v.At + int64(d.SendTime)) < time.Now().Local().Unix() {
						//写入通道，重新推送
						delete(d.Pkgs, k)
						//确定该对象被接收才能解锁
						d.CH <- v
					}
				}
			}
			d.Mx.Unlock()
		case <-d.CloseCH:
			log.Println("数据包管理重发线程退出!")
			close(d.CH)
			return
		}
	}

}

//检测收到的包,如果收到该包返回true
//没有收到数据包返回false
func (d *DataPkgManager) Receive(msgId string) bool {
	d.Mx.RLock()
	if _, ok := d.ReceiveMsgs[msgId]; ok {
		d.Mx.RUnlock()
		return true
	}
	//这个数据包没有收过，需要修改接收状态
	d.Mx.RUnlock()
	//加入接收
	d.ReceiveMsgs[msgId] = true
	//判断是否为发送的数据包，
	//如果是等到发送确认的数据包，修改状态
	d.Mx.Lock()
	_, ok := d.Pkgs[msgId]

	if ok {
		d.Pkgs[msgId].State = true
	}
	d.Mx.Unlock()
	return false
}
func (d *DataPkgManager) AddPkg(msgId, addr, uid string, data []byte) {
	var pkg DataPkg
	pkg.MessageId = msgId
	pkg.State = false
	pkg.Addr = addr
	pkg.Data = data
	pkg.At = time.Now().Local().Unix()
	d.Pkgs[msgId] = &pkg
}

//创建一个消息包管理
func NewDataPkgManager(GCTime, SendTime time.Duration) *DataPkgManager {
	var dpm DataPkgManager
	//默认10秒处理一次GC
	if GCTime == 0 {
		GCTime = 10
	}
	dpm.GCTime = GCTime
	//10秒一次重发检测
	if SendTime == 0 {
		SendTime = 10
	}
	dpm.SendTime = SendTime
	dpm.Pkgs = make(map[string]*DataPkg, 0)
	dpm.CH = make(chan *DataPkg)
	dpm.CloseCH = make(chan int, 1)
	dpm.ReceiveMsgs = make(map[string]bool, 0)
	//启动GC
	go dpm.GC()
	//启动重发检测
	go dpm.ResetSend()
	return &dpm
}

//垃圾数据回收
func (d *DataPkgManager) GC() {
	t := time.NewTicker(d.GCTime * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			d.Mx.Lock()
			for k, v := range d.Pkgs {
				if v.State {
					delete(d.Pkgs, k)
				}
			}
			d.Mx.Unlock()
		case <-d.CloseCH:
			log.Println("数据包管理GC线程已退出!")
			return
		}
	}

}
