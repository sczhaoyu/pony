package main

// import (
// 	"fmt"
// 	"net"
// )

// type Tmo struct {
// 	net.Conn
// 	b bool
// }

// func main() {
// 	var t Tmo
// 	t.b = true
// 	var m map[string]*Tmo = make(map[string]*Tmo)
// 	m["0"] = &t
// 	fmt.Println(m["0"])
// 	//up(&t)
// 	fmt.Println(m["0"])
// }
// func up(t *Tmo) {
// 	t.b = false
// }

// import (
// 	//"fmt"
// 	"sync"
// 	"time"
// )

// func main() {
// 	var m sync.Mutex
// 	ch := make(chan int, 2)
// 	go func() {
// 		for {
// 			ch <- 0
// 			ch <- 1
// 			time.Sleep(time.Second * 5)
// 			m.Lock()
// 		}
// 	}()

// 	go func() {
// 		time.Sleep(time.Second * 2)
// 		m.Lock()
// 		close(ch)

// 		//ch = make(chan int, 20)
// 	}()
// 	// go func() {
// 	// 	for {
// 	// 		ret := <-ch
// 	// 		fmt.Println(ret)

// 	// 	}
// 	// }()
// 	time.Sleep(time.Second * 500)
// }
