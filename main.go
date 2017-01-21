package main

import (
	"fmt"

	"github.com/Alekssasho/GopherInvaders/client"
	"github.com/Alekssasho/GopherInvaders/server"
)

func main() {
	go func() {
		fmt.Println("Server starting")
		server.StartServer()
	}()

	fmt.Println("Client starting")
	client.StartClient()
}
