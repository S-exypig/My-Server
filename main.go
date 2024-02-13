package main

func main() {
	myServer := NewServer("localhost", 8000)
	myServer.Start()
}
