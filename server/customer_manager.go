package server

//客户端多连接管理
type CustomerManager struct {
	Clients map[string]*CustomerServer //所需要链接的服务器
}

//发送消息的路由
func (c *CustomerManager) Router() {

}
