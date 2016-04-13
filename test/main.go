package main

import (
	"fmt"
	"github.com/sczhaoyu/pony/util"
	"net"
)

const (
	addr = "127.0.0.1:8888"
)

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	fmt.Println("已连接服务器")
	defer conn.Close()
	Client(conn)
}

var data []byte = []byte(`{"head":{"command":"100","userId":1976,"token":""}}`)

func Client(conn net.Conn) {
	for {
		tmp := util.IntToByteSlice(len(data))
		tmp = append(tmp, data...)
		conn.Write(tmp)
		buf := make([]byte, 128)
		c, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取服务器数据异常:", err.Error())
			return
		}
		fmt.Println(string(buf[0:c]))
	}

}
