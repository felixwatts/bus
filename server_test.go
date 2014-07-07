package bus

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	var rsp chan (string)

	rsp = make(chan (string))

	_ = newServer()

	time.Sleep(500)

	client, err := net.Dial("tcp", "localhost:6055")
	if err != nil {
		panic(err)
	}

	go read(client, rsp)

	io.WriteString(client, "XX 1 a/b/c\n")

	<-rsp

	io.WriteString(client, "SB 2 cache/i/XEUR/*/price/*\n")

	<-rsp

	io.WriteString(client, "PB 3 cache/i/XEUR/FESX2014090000000/price/bid 121.9\n")

	<-rsp
}

func read(conn net.Conn, rsp chan (string)) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			panic(err)
		}

		s := string(buffer[:n])

		fmt.Print(s)

		rsp <- s
	}
}
