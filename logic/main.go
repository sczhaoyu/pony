package main

import (
	. "github.com/sczhaoyu/pony/logic/server"
)

func main() {
	NewServer(8456).Start()
}
