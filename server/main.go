package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup // threading sync

	// need open 2 port at the same time, can use sync

	// launch file server
	wg.Add(1)
	go startFileServer(&wg)

	// launch Gin Server
	wg.Add(1)
	go startGinServer(&wg)

	// wait 2 server exit
	wg.Wait()
}

// todo
func startGinServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// create gin engine
	r := gin.Default()
	r.LoadHTMLGlob("server/templates/*")

	// add gin router
	// non ui operation
	// 记录回连及其的IP地址

}


func startFileServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// setting file server  15:04
	fs := http.FileServer(http.Dir("."))

	// register file server handler
	http.Handle("/", http.StripPrefix("/", fs))

	// launch file server
	fileServerAddr := "127.0.0.1:8000"

	fmt.Printf("File server is listening on http://%s\n", fileServerAddr)
	err := http.ListenAndServe(fileServerAddr, nil)
	if err != nil {
		fmt.Println("Error starting file server: ", err)
	}
}

// 定义一个全局的map用于存储已写入的IP地址
