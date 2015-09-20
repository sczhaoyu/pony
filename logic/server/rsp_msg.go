package server

type RspMsg struct {
	Rsp      interface{}
	State    bool
	Req      *Request
	SendTime int64
}
