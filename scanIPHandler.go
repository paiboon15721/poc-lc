package main

import (
	"io"
	"net"
	"net/http"
	"strings"

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
		http.Error(w, err.Error(), 500)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).String()
	localIP := localAddr[:strings.IndexByte(localAddr, ':')]
	io.WriteString(w, localIP)
}
