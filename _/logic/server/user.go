package server

func register(c *Conn) {
	c.Out([]byte("111"))
}
func bind(c *Conn) {
	c.Out([]byte("success"))
}
