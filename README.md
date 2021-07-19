# dchat

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/awesome-cmd/dchat?color=pink&logo=go&logoColor=yellow&style=flat-square)
[![Build&test dchat](https://github.com/awesome-cap/im/actions/workflows/Build.yml/badge.svg)](https://github.com/awesome-cap/im/actions/workflows/Build.yml)
[![GitHub license](https://img.shields.io/github/license/awesome-cmd/dchat?color=blue&style=flat-square)](https://github.com/awesome-cap/im/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/awesome-cmd/dchat?color=red&style=flat-square)](https://github.com/awesome-cap/im/stargazers)


# Introduce
**dchat** (Decentralized Chat) 一款去中心化的聊天系统。
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
go get -u github.com/awesome-cap/im
```
## Server:
```powershell
dchat -s -p 3333
```
 - **-s**: 服务端
 - **-p**: 面向客户端的TCP端口，默认为3333
 - **-cluster-port**: 集群通讯端口，默认为3334（可缺省）
 - **-cluster-seeds**: 集群其它部分节点地址，多个用逗号分隔（可缺省）

集群运行示例:
```powershell
dchat -s -p 3333 -cluster-port 3334
dchat -s -p 4001 -cluster-port 4002 -cluster-seeds 127.0.0.1:3334
```
## Client:
```powershell
dchat -c -n nico
```
 - **-c**: 客户端
 - **-n**: 本地昵称
 - **-addrs**: 服务器地址，多个用逗号分隔（可缺省）
 
运行示例:
```powershell
dchat -c -n nico -addrs 127.0.0.1:3333
```
