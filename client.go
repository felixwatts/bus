package bus

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Client interface {
	Stop()
	Subscribe(key string) (string, error)
	Unsubscribe(key string) (string, error)
	Publish(key string, val string) (string, error)
	Rxc() <-chan string
}

type client struct {
	conn net.Conn
	rxc  chan string
	rqId int
}

func (c *client) requestId() string {
	c.rqId++
	return fmt.Sprintf("%v", c.rqId)
}

func (c *client) Rxc() <-chan string {
	return c.rxc
}

func (c *client) Subscribe(keyStr string) (string, error) {
	key, err := parseKey(keyStr)
	if err != nil {
		return "", err
	}

	msg := message{
		meaning:   MSG_TYPE_SUBSCRIBE,
		requestId: c.requestId(),
		key:       key,
	}

	_, err = io.WriteString(c.conn, msg.String())
	if err != nil {
		return "", err
	}

	return msg.requestId, nil
}

func (c *client) Unsubscribe(keyStr string) (string, error) {
	key, err := parseKey(keyStr)
	if err != nil {
		return "", err
	}

	msg := message{
		meaning:   MSG_TYPE_UNSUBSCRIBE,
		requestId: c.requestId(),
		key:       key,
	}

	_, err = io.WriteString(c.conn, msg.String())
	if err != nil {
		return "", err
	}

	return msg.requestId, nil
}

func (c *client) Publish(keyStr string, val string) (string, error) {
	key, err := parseKey(keyStr)
	if err != nil {
		return "", err
	}

	msg := message{
		meaning:   MSG_TYPE_PUBLISH,
		requestId: c.requestId(),
		key:       key,
		val:       val,
	}

	_, err = io.WriteString(c.conn, msg.String())
	if err != nil {
		return "", err
	}

	return msg.requestId, nil
}

func Dial(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &client{
		conn: conn,
		rxc:  make(chan (string), 32),
	}

	go readUntilStop(c)

	return c, nil
}

func readUntilStop(c *client) {
	scanner := bufio.NewScanner(c.conn)
	for {
		more := scanner.Scan()

		if !more {
			break
		}

		msg, err := parseMessage(scanner.Text())

		if err != nil {
			panic(err)
		}

		c.rxc <- msg.String()
	}
}

func (c *client) Stop() {
	c.conn.Close()
}
