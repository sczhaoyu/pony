package server

func register(c *Conn) {
	c.Out("data")
}
func bind(c *Conn) {
	c.Out("success")
}
