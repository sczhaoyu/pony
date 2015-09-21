package server

var ErrMsg map[int]string = make(map[int]string)

func init() {
	initErrMsg()
	registerRoter()
}

func initErrMsg() {
	ErrMsg[-9999] = "token无效"
}
