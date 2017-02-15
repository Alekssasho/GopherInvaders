package main

import (
	"flag"
	"sync"

	"github.com/Alekssasho/GopherInvaders/client"
	"github.com/Alekssasho/GopherInvaders/server"
)

func main() {
	startServer := flag.Bool("server", false, "start game server")
	numberOfPlayer := flag.Int("num-players", 1, "number of players")
	serverIP := flag.String("client-connect", "", "start client connecting to ip")

	flag.Parse()

	var wg sync.WaitGroup

	if *startServer {
		wg.Add(1)
		go func() {
			server.StartServer(*numberOfPlayer)
			wg.Done()
		}()
	}

	if *serverIP != "" {
		// engo cannot run in goroutine so we start in the main goroutine
		client.StartClient(*serverIP)
	}

	//wg.Wait()
}
