# chat
chat in cli.
```go
go get -u github.com/awesome-cmd/chat
```
run server:
```powershell
chat -s -p 3333
```
run client:
```powershell
chat -c -n nico
```
client's usage:
```powershell
[root@centos /]# chat -c nico       // run client
[nico@chat /]# ls                           // get server list
nico                                        
[nico@chat /]# cd nico                      // connect server
[nico@chat nico]# ls                        // get current server's chats
1. nico'schat
[nico@chat nico]# touch golang              // create new chat
Create successful ! The new chat id is 2
[nico@chat nico]# ls                        // get chats again
1. nico'schat
2. golang
[nico@chat nico]# vim 2                     // enter into the chat who id is 2
hello world                                 // enter message what you want to send
nico: hello world           
:q                                          // exit chat
[nico@chat nico]# exit                      // exit client
[root@centos /]# 
```
