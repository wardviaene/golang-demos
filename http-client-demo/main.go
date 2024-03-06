package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main2() {
	resp, err := http.Get("http://localhost:8080/bar")
	if err != nil {
		log.Fatalf("http get failure: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("readall failure: %s\n", err)
	}
	fmt.Printf("body: %s\n", body)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("connection failure: %s\n", err)
	}
	fmt.Fprintf(conn, "GET /bar HTTP/1.0\r\nHost: localhost\r\n\r\n")
	buf := bufio.NewReader(conn)
	status, err := buf.ReadString('\n')
	if err != nil {
		log.Fatalf("readstring error: %s\n", err)
	}
	fmt.Printf("status: %s", status)
	for {
		line, err := buf.ReadString('\n')
		fmt.Printf("%s", line)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("readstring error: %s\n", err)
		}
	}
	err = conn.Close()
	if err != nil {
		log.Fatalf("connection close error: %s\n", err)
	}
}
