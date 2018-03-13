package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func scanIPHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Get localIP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).String()
	localIP := localAddr[:strings.IndexByte(localAddr, ':')]

	// Scan /24 subnet
	incIP := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}
	ip, ipNet, _ := net.ParseCIDR(fmt.Sprintf("%s/24", localIP))
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		fmt.Println(ip.String())
	}
	io.WriteString(w, localIP)
}
