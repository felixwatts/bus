package bus

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type clientHandler struct {
	id     int
	server *server
	conn   net.Conn
	sendc  chan (string)
	recvc  chan (string)
	closec chan (bool)
}

func newClientHandler(conn net.Conn, id int, server *server) *clientHandler {
	clientHandler := &clientHandler{
		id:     id,
		server: server,
		conn:   conn,
		sendc:  make(chan (string), 32),
		recvc:  make(chan (string), 32),
		closec: make(chan (bool)),
	}

	go clientHandler.run()
	go clientHandler.listen()

	return clientHandler
}

func (clientHandler *clientHandler) run() {

	defer clientHandler.close()

	for {
		select {
		case message := <-clientHandler.sendc:
			clientHandler.log(fmt.Sprintf("TX %v", message))
			_, err := io.WriteString(clientHandler.conn, message)
			if err != nil {
				clientHandler.log(err)
				return
			}
			break
		case <-clientHandler.closec:
			return
		}
	}

}

func (clientHandler *clientHandler) send(message string) {
	clientHandler.sendc <- message
}

func (clientHandler *clientHandler) sendOK(requestId string) {
	clientHandler.send(message{meaning: MSG_TYPE_OK, requestId: requestId}.String())
}

func (clientHandler *clientHandler) sendFail(requestId string) {
	clientHandler.send(message{meaning: MSG_TYPE_FAIL, requestId: requestId}.String())
}

func (clientHandler *clientHandler) close() {
	clientHandler.log("Closing")
	clientHandler.server.removeClient(clientHandler)
	clientHandler.conn.Close()
}

func (clientHandler *clientHandler) listen() {

	defer clientHandler.log("End listen")

	clientHandler.log("Listening")

	scanner := bufio.NewScanner(clientHandler.conn)

	for {
		more := scanner.Scan()

		if !more {
			break
		}

		clientHandler.log(fmt.Sprintf("RX [%v]", scanner.Text()))

		request := request{
			scanner.Text(),
			clientHandler,
		}

		clientHandler.server.requestc <- request
	}
}

func (clientHandler *clientHandler) log(data interface{}) {
	log.Printf("[CLIENT #%v] %v", clientHandler.id, data)
}
