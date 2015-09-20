package server

type RequestHeader struct {
	UserId   int64  `json:"userId"`
	Addr     string `json:"addr"`
	FaceCode int    `json:"faceCode"`
	Token    string `json:"token"`
}

type Request struct {
	Head *RequestHeader `json:"head"`
	Body []byte         `json:"body"`
}
