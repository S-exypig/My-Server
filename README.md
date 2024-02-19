# My-Server
使用go语言搭建自己的一个网络服务器，学习向

## 2/19
简易服务器已搭建完成，新增客户端的部分实现

## 2/20
客户端已全部实现

## 使用说明
clone项目到本地后，进入对应文件夹，打开shell键入：
```sh
go run ./server.go ./user.go ./main.go -ip SERVER_IP -port SERVER_PORT
```
以运行服务端程序，服务端默认IP和Port为localhost和8000。

之后在新的shell中键入:
```sh
go run ./client.go -ip SERVER_IP -port SERVER_PORT
```
以运行客户端程序，服务端默认IP和Port为localhost和8000。