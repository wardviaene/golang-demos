package main

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
)

func main2() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen error: %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("accept error: %s", err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	fmt.Printf("client is connected: %s\n", conn.RemoteAddr())

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Buffer: %s\n", buf)

	_, err = conn.Write([]byte("HTTP/1.0 200 OK\nDate: Mon, 04 Mar 2024 16:18:16 GMT\nContent-Length: 13\nContent-Type: text/plain; charset=utf-8\n\nHello, \"/bar\""))
	if err != nil {
		log.Fatalf("write error: %s", err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatalf("write error: %s", err)
	}
}
