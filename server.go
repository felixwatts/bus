package bus

import (
	"fmt"
	"log"
	"net"
)

type request struct {
	message string
	client  *client
}

type server struct {
	nextClientId  int
	ln            net.Listener
	clients       map[int]*client
	requestc      chan (request)
	closec        chan (bool)
	subscriptions *keyTree
}

func newServer() *server {
	server := &server{
		clients:       make(map[int]*client),
		requestc:      make(chan (request), 1024),
		closec:        make(chan (bool)),
		subscriptions: newKeyTree(),
	}

	go server.run()
	go server.serve()

	return server
}

func (server *server) serve() {

	ln, err := net.Listen("tcp", "localhost:6055")
	if err != nil {
		server.log(err)
		return
	}
	server.ln = ln

	server.log("Ready")

	for {
		conn, err := ln.Accept()
		if err != nil {
			server.log(err)
			return
		}

		client := newClient(conn, server.nextClientId, server)

		server.clients[client.id] = client

		server.nextClientId++
	}
}

func (server *server) run() {
	for {
		select {
		case request := <-server.requestc:
			server.handleRequest(request)
			break
		case <-server.closec:
			server.ln.Close()
			return
		}
	}
}

func (server *server) handleRequest(request request) {

	server.log(fmt.Sprintf("Handle %v", request.message))

	message, err := parseMessage(request.message)
	if err != nil {
		server.logBadRequest(request, err)
		return
	}

	switch message.meaning {
	case MSG_TYPE_SUBSCRIBE:
		server.subscriptions.subscribe(request.client, message.key)
		request.client.sendOK(message.requestId)
		break
	case MSG_TYPE_PUBLISH:
		server.subscriptions.publish(message.key, message.String(), false)
		request.client.sendOK(message.requestId)
		break
	default:
		server.log(fmt.Sprintf("Unknown message type [%v] from #%v", message.meaning, request.client.id))
		request.client.sendFail(message.requestId)
	}
}

func (server *server) close() {
	server.log("Closing")
	close(server.closec)
}

func (server *server) removeClient(client *client) {

}

func (server *server) log(data interface{}) {
	log.Printf("[SERVER] %v\n", data)
}

func (server *server) logBadRequest(request request, err error) {
	server.log(fmt.Sprintf("Bad request from #%v (%v) [%v]", request.client.id, request.message, err))
}
