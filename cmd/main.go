package main

import (
	"fmt"
	"github.com/Ladence/go-url-shortener/internal/server"
	"github.com/Ladence/go-url-shortener/internal/storage"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixMilli())
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
	fmt.Println("storages are ready, going to start server")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGKILL, syscall.SIGINT)
	appServer := server.NewServer("127.0.0.1:7002", urlStorage, ipStorage)
	go appServer.Run()
	sig := <-sigChan
	fmt.Printf("got signal: %s. closing an application server", sig.String())
	err = appServer.Shutdown()
	if err != nil {
		fmt.Printf("error on shutdown application server: %v", err)
	}
}
