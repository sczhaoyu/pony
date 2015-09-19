package server

type requestHead struct {
	UserId   int64
	FaceCode int
	Uuid     string
}
type request struct {
	Head *requestHead
}
