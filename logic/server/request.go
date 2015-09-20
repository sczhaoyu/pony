package server

type RequestHead struct {
	UserId   int64
	FaceCode int
	Uuid     string
}
type Request struct {
	Head *RequestHead
}
