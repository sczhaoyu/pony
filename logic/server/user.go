package server

func register(r *Request) {
	s.pushSendQueue(r, "注册完成", false)
}
