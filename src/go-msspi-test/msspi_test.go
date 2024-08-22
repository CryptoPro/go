package msspitest

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"
)

func clientGet(t *testing.T, host string, uri string) []byte {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, 44333), 2*time.Second)
	if err != nil {
		t.Fatalf("Unexpected error on dial: %v", err)
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
	defer tlsConn.Close()

	wbuf := []byte("GET " + uri + " HTTP/1.1\r\nHost: " + host + "\r\n\r\n")
	if wlen, err := tlsConn.Write(wbuf); wlen != len(wbuf) || err != nil {
		t.Fatalf("Error sending: %v", err)
	}

	rbuf := make([]byte, 16384)
	if rlen, err := tlsConn.Read(rbuf); rlen == 0 || err != nil {
		t.Fatalf("Error reading: %v", err)
		return nil
	} else {
		return rbuf[:rlen]
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func RunServer() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":44333", "server.crt", "server.crt", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func TestMsspiServer(t *testing.T) {
	go RunServer()
	rbuf := clientGet(t, "localhost", "/hello")
	s := string(rbuf)
	fmt.Printf(s + "\n")
	for i := 0; i < 9999999; i++ {
		time.Sleep(time.Second)
		runtime.GC()
	}
}
