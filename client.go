package bus

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type client struct {
	id     int
	server *server
	conn   net.Conn
	sendc  chan (string)
	recvc  chan (string)
	closec chan (bool)
}

func newClient(conn net.Conn, id int, server *server) *client {
	client := &client{
		id:     id,
		server: server,
		conn:   conn,
		sendc:  make(chan (string), 32),
		recvc:  make(chan (string), 32),
		closec: make(chan (bool)),
	}

	go client.run()
	go client.listen()

	return client
}

func (client *client) run() {

	defer client.close()

	for {
		select {
		case message := <-client.sendc:
			client.log(fmt.Sprintf("TX %v", message))
			_, err := io.WriteString(client.conn, message)
			if err != nil {
				client.log(err)
				return
			}
			break
		case <-client.closec:
			return
		}
	}

}

func (client *client) send(message string) {
	client.sendc <- message
}

func (client *client) sendOK(requestId string) {
	client.send(message{meaning: MSG_TYPE_OK, requestId: requestId}.String())
}

func (client *client) sendFail(requestId string) {
	client.send(message{meaning: MSG_TYPE_FAIL, requestId: requestId}.String())
}

func (client *client) close() {
	client.log("Closing")
	client.server.removeClient(client)
	client.conn.Close()
}

func (client *client) listen() {

	defer client.log("End listen")

	client.log("Listening")

	scanner := bufio.NewScanner(client.conn)

	for {
		more := scanner.Scan()

		if !more {
			break
		}

		client.log(fmt.Sprintf("RX [%v]", scanner.Text()))

		request := request{
			scanner.Text(),
			client,
		}

		client.server.requestc <- request
	}
}

func (client *client) log(data interface{}) {
	log.Printf("[CLIENT #%v] %v", client.id, data)
}
