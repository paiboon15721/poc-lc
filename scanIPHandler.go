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
		wg  sync.WaitGroup
		mux sync.Mutex
	)
	type osInfo struct {
		DistributorID string `json:"distributorID,omitempty"`
		Description   string `json:"description,omitempty"`
		Release       string `json:"release,omitempty"`
		Codename      string `json:"codename,omitempty"`
	}
	type quota struct {
		Total  int `json:"total,omitempty"`
		Remain int `json:"remain,omitempty"`
	}
	type firmwareInfo struct {
		Version    string `json:"version,omitempty"`
		BuildTime  string `json:"buildTime,omitempty"`
		HardwareID string `json:"hardwareID,omitempty"`
		Customer   string `json:"customer,omitempty"`
		Quota      quota  `json:"quota,omitempty"`
	}
	type serverInfo struct {
		IP           string        `json:"ip"`
		OsInfo       *osInfo       `json:"osInfo,omitempty"`
		FirmwareInfo *firmwareInfo `json:"firmwareInfo,omitempty"`
	}
	var serverInfos []serverInfo

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
				currentServerInfo := serverInfo{IP: ip}
				mux.Lock()
				serverInfos = append(serverInfos, currentServerInfo)
				mux.Unlock()
			}
		}(ip.String())
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverInfos)
}
