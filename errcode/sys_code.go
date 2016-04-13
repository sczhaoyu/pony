package errcode

type Error struct {
	Code int    //错误代码
	Msg  string //错误消息
	error
}

var (
	CLIENT_TIME_OUT = NewError(100, "人员已满请稍后再试！")
	READ_TIME_OUT   = NewError(101, "网络传输超时！")
	READ_DATA_MAX   = NewError(102, "数据传输已超出最大")
	READ_DATA_NULL  = NewError(103, "数据不能为空")
)

func (e *Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *Error {
	var ret Error
	ret.Code = code
	ret.Msg = msg
	return &ret
}
func RestErr(err *Error, msg string) *Error {
	var em Error
	em.Code = err.Code
	em.Msg = msg
	return &em
}
