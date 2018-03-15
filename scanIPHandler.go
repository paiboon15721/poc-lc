package main

import (
	"encoding/json"
	"fmt"
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
	// var wg sync.WaitGroup
	var detectedIPs []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		detectedIPs = append(detectedIPs, ip.String())
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detectedIPs)
}
