package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	ssdpAddr    = "239.255.255.250:1982"
	discoverMsg = "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1982\r\n MAN:\"ssdp:discover\"\r\n ST:wifi_bulb\r\n"
	crlf        = "\r\n"
)

// Smartbulb represents a instance of smartbulb
type Smartbulb struct {
	addr string
}

func main() {
	if err := discover(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func discover() error {
	ssdp, _ := net.ResolveUDPAddr("udp4", ssdpAddr)
	c, _ := net.ListenPacket("udp4", ":0")

	socket := c.(*net.UDPConn)
	socket.WriteToUDP([]byte(discoverMsg), ssdp)
	socket.SetReadDeadline(time.Now().Add(3 * time.Second))

	rsBuf := make([]byte, 1024)

	size, _, err := socket.ReadFromUDP(rsBuf)

	if err != nil {
		return errors.New("No Devices Found")
	}

	rs := rsBuf[0:size]
	addr := parseAddr(string(rs))
	fmt.Printf("Device with IP %s found\n", addr)
	return nil
}

func parseAddr(msg string) string {
	if strings.HasSuffix(msg, crlf) {
		msg = msg + crlf
	}
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(msg)), nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	return strings.TrimPrefix(resp.Header.Get("LOCATION"), "yeelight://")
}
