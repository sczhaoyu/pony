package server

func register(r *request) {
	s.pushSendQueue(r, "注册完成", false)
}
