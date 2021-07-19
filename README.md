# Introduce
去中心化的聊天系统。
# Features
 - 轻量级
 - Unix指令交互
 - 去中心化
 - 断线重连
 - 支持集群
 - 分布式ID
# Start
## Install
```golang
go get -u github.com/awesome-cap/im/cmd/imc     //获取客户端
go get -u github.com/awesome-cap/im/cmd/ims     //获取服务端
```
## Server:
```powershell
ims -p 3333
```
 - **-p**: 面向客户端的TCP端口，默认为3333
 - **-cluster-port**: 集群通讯端口，默认为3334（可缺省）
 - **-cluster-seeds**: 集群其它部分节点地址，多个用逗号分隔（可缺省）

集群运行示例:
```powershell
ims -p 3333 -cluster-port 3334
ims -p 4001 -cluster-port 4002 -cluster-seeds 127.0.0.1:3334
```
## Client:
```powershell
imc -n nico
```
 - **-n**: 本地昵称
 - **-addrs**: 服务器地址，多个用逗号分隔（可缺省）
 
运行示例:
```powershell
imc -c -n nico -addrs 127.0.0.1:3333
```
