package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

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

	// Decare datastructure to hold server information
	var (
		wg          sync.WaitGroup
		mux         sync.Mutex
		detectedIPs []string
	)
	type osInfo struct {
		DistributorID string `json:"distributorID"`
		Description   string `json:"description"`
		Release       string `json:"release"`
		Codename      string `json:"codename"`
	}
	type quota struct {
		Total  int `json:"total"`
		Remain int `json:"remain"`
	}
	type firmwareInfo struct {
		Version    string `json:"version"`
		BuildTime  string `json:"buildTime"`
		HardwareID string `json:"hardwareID"`
		Customer   string `json:"customer"`
		Quota      quota  `json:"quota"`
	}
	type serverInfo struct {
		IP           string       `json:"ip"`
		OsInfo       osInfo       `json:"osInfo"`
		FirmwareInfo firmwareInfo `json:"firmwareInfo"`
	}

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
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:22", ip), time.Millisecond)
			if err == nil {
				c.Close()
				mux.Lock()
				detectedIPs = append(detectedIPs, ip)
				mux.Unlock()
			}
		}(ip.String())
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detectedIPs)
}
