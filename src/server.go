package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"log"
	"crypto/sha1"
	"encoding/hex"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErrorServer(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErrorServer(err)
	f, err := os.OpenFile("../../log/log.txt", os.O_RDWR | os.O_CREATE | os.O_CREATE, 0666)
	checkErrorServer(err)
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
	log.SetOutput(f)

	for {
		conn, err := listener.Accept()
		checkErrorServer(err)
		conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		now := []byte(string(strconv.FormatInt(time.Now().UnixNano(), 10)))
		hash := sha1.New()
		hash.Write(now)
		id := hex.EncodeToString(hash.Sum(nil))
		go handleClient(conn, id)
	}

	defer f.Close()
}

func handleClient(conn net.Conn, id string) {
	defer conn.Close()
	for {
		request := make([]byte, 128)
		read_len, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}

		log.Println(id + " request => " + string(request[0:read_len]))

		if string(request[0:read_len]) == "timestamp" {
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			conn.Write([]byte(timestamp))
		} else if string(request[0:read_len]) == "echo" {
			echo := string(request[0:read_len])
			conn.Write([]byte(echo))
		} else {
			daytime := time.Now().String()
			conn.Write([]byte(daytime))
		}
	}
}

func checkErrorServer(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
