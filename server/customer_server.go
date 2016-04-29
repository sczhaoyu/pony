package server

import (
	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/sczhaoyu/pony/util"
	"log"
	"net"
	"time"
)

//该类主要用于客户端链接
//接收和回应数据包，前4字节为数据包大小
type CustomerServer struct {
	net.Conn
	Name           string                         //服务器名称
	Id             string                         //身份标示ID
	State          bool                           //链接状态
	ResetChan      chan int                       //重置通道信号
	ResetTimeOut   int                            //超时重链接秒
	DataChan       chan []byte                    //数据发送通道
	MaxData        int                            //最大读取限制
	ServerAddr     string                         //服务器链接地址 IP+Port格式
	Handler        func(*CustomerServer, *Respon) //数据读取处理
	FirstSend      func()                         //启动或者重新链接第一次发送消息
	HeartbeatTime  int64                          //心跳时间秒
	DPM            *DataPkgManager                //数据包重发管理
	Body           []byte                         //当前回应的body数据包
	CloseState     bool                           //关闭状态
	CloseHeartbeat chan int                       //心跳关闭通道
}

//关闭链接，销毁
func (c *CustomerServer) CloseClient() {
	//设置状态为关闭
	c.CloseState = true
	c.Conn.Close()
	//关闭数据通道
	close(c.DataChan)
	//关闭心跳
	c.CloseHeartbeat <- 1
	//数据包关闭
	close(c.DPM.CloseCH)

}
func (c *CustomerServer) Unmarshal(b interface{}) error {
	if len(c.Body) == 0 {
		return errors.New("body data nil!")
	}
	return json.Unmarshal(c.Body, b)
}

//循环发送心跳
func (c *CustomerServer) heartbeat() {
	t := time.NewTicker(time.Second * 3)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if c.State {
				//发送心跳
				c.WriteJson(nil, 520)
			}
		case <-c.CloseHeartbeat:
			log.Println("心跳线程已退出!")
			return
		}
	}

}

//数据包重发
func (c *CustomerServer) rewire() {
	for {
		pkg, isClose := <-c.DPM.CH
		if !isClose {
			log.Println("重发线程已退出!")
			return
		}
		var req Request
		req.Unmarshal(pkg.Data[0:4])
		//往服务器推送
		c.Write(pkg.Data)
		//加入数据包检测
		c.DPM.AddPkg(req.Header.RequestId, c.RemoteAddr().String(), "", pkg.Data)
	}
}

//发送Json
func (c *CustomerServer) WriteJson(b interface{}, faceCode int) {

	var req Request
	//请求数据的ID
	req.Header.RequestId = util.GetUUID()
	req.Header.FaceCode = faceCode
	req.Header.RequestTime = time.Now().Unix()
	req.Header.UserAddr = c.LocalAddr().String()
	req.Body = b
	d := util.ByteLen(req.Marshal())
	//加入重发检测
	if faceCode != 520 && faceCode != 200 {
		c.DPM.AddPkg(req.Header.RequestId, c.RemoteAddr().String(), "", d)
	}
	//写入通道
	c.DataChan <- d

}

//创建连接
func (c *CustomerServer) init() {
	//数据包重发管理
	c.CloseState = false
	c.DPM = NewDataPkgManager(5, 5)
	if c.HeartbeatTime <= 0 {
		c.HeartbeatTime = 5
	}
	if c.Name == "" {
		c.Name = c.ServerAddr
	}
	if c.Id == "" {
		c.Id = util.GetUUID()
	}
	c.Conn = new(Conn)
	//默认重连接20秒
	if c.ResetTimeOut == 0 {
		c.ResetTimeOut = 20
	}
	c.ResetChan = make(chan int, 1)
	c.CloseHeartbeat = make(chan int, 1)
	if len(c.DataChan) == 0 {
		c.DataChan = make(chan []byte, 100)
	}
	if c.MaxData == 0 {
		c.MaxData = 10000000
	}
}
func (c *CustomerServer) Run() {
	c.init()
	go func() {
		conn, err := c.NewConn()
		if err != nil {
			c.State = false
			c.ResetChan <- 0
		} else {
			c.Conn = conn
			c.State = true
			log.Println(fmt.Sprintf("[%s]start success addr:%s", c.Name, c.ServerAddr))
			if c.FirstSend != nil {
				c.FirstSend()
			}
			go c.Read()
		}

	}()
	//链接状态监测
	go c.CheckClient()
	//数据发送线程
	go c.SendData()
	//启动心跳
	go c.heartbeat()
	//数据重发
	go c.rewire()

}

//读取数据
func (c *CustomerServer) Read() {
	for {
		data, err := util.ReadData(c.Conn, c.MaxData)
		//数据读取出错，不处理
		if err != nil {
			c.State = false
			c.ResetChan <- 0
			break
		}
		//判断处理器是否实现
		if c.Handler != nil {
			var rsp Respon
			err = rsp.Unmarshal(data)
			//如果反序错误不处理
			if err != nil {
				log.Println("server respons err:", err.Error())
			} else {
				//判断该数据包是否接收过，如果接收过不需要处理
				if ok := c.DPM.Receive(rsp.Header.ResponsId); ok {
					//跳过该数据包的处理
					continue
				}
				c.DPM.Receive(rsp.Header.RequestId)
				//回应服务器已经收包
				c.WriteJson(rsp.Header.ResponsId, 200)
				//取出data里面的body
				j, _ := simplejson.NewJson(data)
				body, err := j.Get("body").MarshalJSON()
				if err == nil {
					c.Body = body
				}
				c.Handler(c, &rsp)
			}

		}
	}
}

//阻塞发送数据
func (c *CustomerServer) SendData() {
	for {
		data, isClose := <-c.DataChan
		if !isClose {
			log.Println("发送线程已退出!")
			return
		}
		c.Conn.Write(data)
	}

}

//检查链接的完整
func (c *CustomerServer) CheckClient() {
	for {
		<-c.ResetChan
		for c.State == false {
			//说明是关闭状态没直接返回
			if c.CloseState {
				log.Println("重连线程已退出!")
				return
			}
			var err error
			c.Conn, err = c.NewConn()
			if err == nil {
				c.State = true
				if c.FirstSend != nil {
					c.FirstSend()
				}
				log.Println(fmt.Sprintf("[%s]reset success addr:%s", c.Name, c.ServerAddr))
				go c.Read()
			} else {
				log.Println(fmt.Sprintf("[%s]waiting client %s fail!", c.Name, c.ServerAddr))
				time.Sleep(time.Second * time.Duration(c.ResetTimeOut))
			}
		}
	}
}

//创建一个链接
func (c *CustomerServer) NewConn() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
