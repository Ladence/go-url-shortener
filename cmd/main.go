package main

import (
	"fmt"
	"github.com/Ladence/go-url-shortener/internal/server"
	"github.com/Ladence/go-url-shortener/internal/storage"
)

func main() {
	urlStorage, err := storage.NewRedisStorage("127.0.0.1:7005", "", 0)
	if err != nil {
		fmt.Printf("error on initializing urlStorage: %v", err)
		return
	}
	ipStorage, err := storage.NewRedisStorage("127.0.0.1:7005", "", 1)
	if err != nil {
		fmt.Printf("error on initializing urlStorage: %v", err)
		return
	}
	fmt.Println("Storages are ready, going to start server")
	appServer := server.NewServer("127.0.0.1:7002", urlStorage, ipStorage)
	appServer.Run()
}
