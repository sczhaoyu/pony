admin这个包负责管理所有的业务服务器和查看每台服务器的状态

business 这个包负责处理client发送过来的任务，处理具体的业务逻辑


client 包负责与客户端保持连接



关于启动  


请以此启动  admin、business、 client 三个包下的main.go

test包就是一个测试的东西，发送命令什么的。

启动test.go 就可以看到
