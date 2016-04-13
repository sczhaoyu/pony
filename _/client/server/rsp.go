package server

type Rsp struct {
	Data      []byte //回应的数据
	SessionId string //通知的用户
}
