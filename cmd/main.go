package main

import (
	"fmt"
	"log"

	"github.com/KyberNetwork/node-monitor/server"
)

func main() {
	fmt.Println("hello world!")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := server.NewServer()
	server.Run()
}
