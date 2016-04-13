package customer

import (
	. "github.com/sczhaoyu/pony/server"
	"log"
)

var router map[int]func(*CustomerServer, *Respon) = make(map[int]func(*CustomerServer, *Respon))

func init() {
	router[100] = login
	router[200] = ok
}
func Handle(c *CustomerServer, rsp *Respon) {
	if fun, ok := router[rsp.Header.FaceCode]; ok {
		fun(c, rsp)
	}
}
func login(c *CustomerServer, r *Respon) {
	log.Println("login success!")
}
func ok(c *CustomerServer, r *Respon) {
	log.Println("received  request id:", r.Header.RequestId)
}
