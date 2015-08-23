package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port <mode>", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErrorClient(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkErrorClient(err)
	result := make([]byte, 128)
	mode := "date"
	if len(os.Args) > 2{
		mode = os.Args[2]
	}

	for {
		_, err := conn.Write([]byte(mode))
		checkErrorClient(err)
		conn.Read(result)
		checkErrorClient(err)
		fmt.Println(string(result))
		time.Sleep(1 * time.Second)
	}
	os.Exit(0)
}

func checkErrorClient(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
