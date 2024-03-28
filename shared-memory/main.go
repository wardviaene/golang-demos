package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", indexHandler)

	var wg sync.WaitGroup

	result := make([]string, 0, 1)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			fmt.Printf("Start goroutine %d\n", i)
			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			result[i] = string(body)
			//result = append(result, string(body))

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		fmt.Printf("%+v", result)
	}()

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		log.Fatalf("random generator error: %s", err)
	}
	str := strconv.FormatInt(nBig.Int64(), 10)
	time.Sleep(5 * time.Second)
	w.Write([]byte(str))
}
