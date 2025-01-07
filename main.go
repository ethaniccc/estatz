package main

import "github.com/ethaniccc/estatz/server"

func main() {
	server.New(server.Config{
		Port: 10000,
	})
}
