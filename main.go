package main

import "flag"

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "localhost", "指定服务端的ip地址")
	flag.IntVar(&serverPort, "port", 8000, "指定服务端的端口号")
}

func main() {
	flag.Parse()
	myServer := NewServer(serverIp, serverPort)
	myServer.Start()
}
