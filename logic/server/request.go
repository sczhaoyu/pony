package server

type RequestHead struct {
	UserId   int64  `json:"userId"`
	FaceCode int    `json:"faceCode"`
	Uuid     string `json:"uuid"`
}
type Request struct {
	Head *RequestHead `json:"head"`
}
