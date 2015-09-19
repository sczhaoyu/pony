package server

type ResponseHead struct {
	UserId int64  `json:"userId"`
	Uuid   string `json:"uuid"`
}
type Response struct {
	Head *ResponseHead `json:"head"`
}

func NewClientResponse(r *request, data interface{}) *Response {
	return nil
}
func (r *Response) GetJson() []byte {
	return nil
}
