package server

type RspMsg struct {
	Rsp      interface{}
	State    bool
	Req      *request
	SendTime int64
}
