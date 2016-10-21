package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const (
	host     = "0.0.0.0"
	port     = "18080"
	respFile = "resp.json"
)

// Resp is struct of response of UDP server
type Resp struct {
	Instances int      `json:"instances"`
	ConnNum   int      `json:"connNum"`
	OvsID     int      `json:"ovsId"`
	ScaleInIP string   `json:"ScaleInIp"`
	LiveGWs   []string `json:"LiveGWs"`
}

func main() {
	fmt.Println("started...")
	addr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()
	for {
		handleClient(conn)
	}
}
func handleClient(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}
	fmt.Println(n, remoteAddr)
	fmt.Printf("Req: %s\n", string(data))

	resp, _ := getResp()

	respData, _ := json.Marshal(resp)

	fmt.Printf("Resp: %s\n", string(respData))

	conn.WriteToUDP([]byte(respData), remoteAddr)
}

func getResp() (resp *Resp, err error) {
	content, err := readTextFile(respFile)
	if err != nil {
		log.Printf("read file error: %v\n", err)
		return
	}
	resp, err = parseJSON(content)
	if err != nil {
		log.Printf("parse json error: %v\n", err)
		return
	}
	return
}

func parseJSON(content []byte) (resp *Resp, err error) {
	resp = &Resp{}
	err = json.Unmarshal(content, resp)
	if err != nil {
		log.Printf("unmarshal to json error: %v\n", err)
		return
	}
	return
}

func readTextFile(path string) (content []byte, err error) {
	if _, err = os.Stat(path); err != nil {
		log.Printf("stat file error: %v\n", err)
		return
	}
	content, err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read file error: %v\n", err)
		return
	}
	return content, nil
}
