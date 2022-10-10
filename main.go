package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
)

func main() {
	port := flag.Int("port", 4001, "Listen port")
	apiUrl := flag.String("api", "http://127.0.0.1/api/v2/tracks/ais", "Send URL")

	serve(*port, *apiUrl)
}

func serve(port int, apiUrl string) {
	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: port, Zone: ""})
	if err != nil {
		panic("Server stopped.")
	}

	defer func(ServerConn *net.UDPConn) {
		err := ServerConn.Close()
		if err != nil {
			panic("Stopped on close.")
		}
	}(ServerConn)

	buf := make([]byte, 1024)

	for {
		n, _, _ := ServerConn.ReadFromUDP(buf)
		sendToApi(string(buf[0:n]), apiUrl)
	}
}
func sendToApi(msg string, apiUrl string) {
	url := apiUrl
	values := map[string]string{"msg": msg}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Print("Cannot send.")
	}

	if err == nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Print("Cannot close.")
			}
		}(resp.Body)
	}
}
