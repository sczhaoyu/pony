package server

func handler(r *RoterServer, data []byte) {
	var w Response
	err := w.Unmarshal(data)
	if err != nil {
		return
	}
	//建立此基础的必要条件 必须和连接端服务session同步
	// ONE  = 1 //单个通知
	// More = 2 //多个用户通知
	// All  = 3 //全局通知
	//读取Users的数据节点
	//在session中查找他们所在的服务器
	switch w.Head.NoticeType {
	case ONE:
		//找到单个用户，通知
	case More:
		//取出一部分的用户ID通知
	case All:
		//无需判断通知全部客户端的链接
	}
}
