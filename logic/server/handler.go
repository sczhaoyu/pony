package server

import (
	"encoding/json"
)

func handler(c *Conn, data []byte) {
	err := json.Unmarshal(data, &c.Request)
	if err != nil {
		c.Out([]byte(err.Error()))
		return
	}
	b := beforeInterceptor
	for i := 0; i < len(b); i++ {
		if b[i](c) == false {
			return
		}
	}
	if ReadFunc != nil {
		ReadFunc(c)
	}
	a := afterInterceptor
	for i := 0; i < len(a); i++ {
		if a[i](c) == false {
			return
		}
	}
}
