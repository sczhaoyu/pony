package main

import (
	admin "github.com/sczhaoyu/pony/admin/admin_server"
)

func main() {
	admin.NewAdminServer(2058).Run()
}
