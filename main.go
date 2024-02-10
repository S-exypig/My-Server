package main

import (
	"myserver/server"
)

func main() {
	myServer := server.NewServer("localhost", 8000)
	myServer.Start()
}
