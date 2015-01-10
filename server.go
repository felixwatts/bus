package bus

import (
	"fmt"
	"log"
	"net"
)

type Server interface {
	Stop()
}

type request struct {
	message string
	client  *clientHandler
}

type server struct {
	nextClientId  int
	ln            net.Listener
	clients       map[int]*clientHandler
	requestc      chan (request)
	closec        chan (bool)
	subscriptions *keyTree
}

func Serve(addr string) (Server, error) {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := &server{
		clients:       make(map[int]*clientHandler),
		requestc:      make(chan (request), 1024),
		closec:        make(chan (bool)),
		subscriptions: newKeyTree(),
		ln:            ln,
	}

	go server.run()
	go server.serve()

	return server, nil
}

func (server *server) Stop() {
	close(server.closec)
}

func (server *server) serve() {

	server.log(fmt.Sprintf("Listening on %v", server.ln.Addr()))

	for {
		conn, err := server.ln.Accept()
		if err != nil {
			server.log(err)
			return
		}

		client := newClientHandler(conn, server.nextClientId, server)

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
	case MSG_TYPE_UNSUBSCRIBE:
		server.subscriptions.unsubscribe(request.client, message.key)
		request.client.sendOK(message.requestId)
		break
	case MSG_TYPE_CLAIM:

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

func (server *server) removeClient(client *clientHandler) {

}

func (server *server) log(data interface{}) {
	log.Printf("[SERVER] %v\n", data)
}

func (server *server) logBadRequest(request request, err error) {
	server.log(fmt.Sprintf("Bad request from #%v (%v) [%v]", request.client.id, request.message, err))
}
