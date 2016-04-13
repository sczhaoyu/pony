package server

var (
	beforeInterceptor []func(*Conn) bool
	afterInterceptor  []func(*Conn) bool
)

func BeforeInterceptor(f func(*Conn) bool) {
	beforeInterceptor = append(beforeInterceptor, f)
}
func AfterInterceptor(f func(*Conn) bool) {
	afterInterceptor = append(afterInterceptor, f)
}
