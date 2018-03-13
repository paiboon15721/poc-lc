package main

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// func incIP(ip net.IP) {
// 	for j := len(ip) - 1; j >= 0; j-- {
// 		ip[j]++
// 		if ip[j] > 0 {
// 			break
// 		}
// 	}
// }

func scanIPHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	io.WriteString(w, localAddr.String())
}
