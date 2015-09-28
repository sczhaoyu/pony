package server

var (
	ReadFunc func(*Conn)
)

func init() {
	steupSysRoter()
}
