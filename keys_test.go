package bus

import (
	"testing"
	"time"
)

var subject *keyTree

func Test1(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key, _ := parseKey("a")

	subject.subscribe(client1, key)
	subject.publish(key, "msg", false)

	client1.assertReceives("msg", t)
}

func Test2(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)
	client2 := mockClient(2)

	key, _ := parseKey("a")

	subject.subscribe(client1, key)
	subject.subscribe(client2, key)
	subject.publish(key, "msg", false)

	client1.assertReceives("msg", t)
	client2.assertReceives("msg", t)
}

func Test3(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)
	client2 := mockClient(2)

	key1, _ := parseKey("a")
	key2, _ := parseKey("b")

	subject.subscribe(client1, key1)
	subject.subscribe(client2, key2)
	subject.publish(key1, "msg", false)

	client1.assertReceives("msg", t)
	client2.assertReceivesNothing(t)
}

func Test4(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key, _ := parseKey("a/b")

	subject.subscribe(client1, key)
	subject.publish(key, "msg", false)

	client1.assertReceives("msg", t)
}

func Test5(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("a/*")
	key2, _ := parseKey("a/b")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceives("msg", t)
}

func Test6(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("*/b")
	key2, _ := parseKey("a/b")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceives("msg", t)
}

func Test7(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("*/c")
	key2, _ := parseKey("a/b")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceivesNothing(t)
}

func Test8(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("c/b")
	key2, _ := parseKey("a/b")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceivesNothing(t)
}

func Test9(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("a/**")
	key2, _ := parseKey("a/b/c")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceives("msg", t)
}

func Test10(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("a/**/d")
	key2, _ := parseKey("a/b/c/d")

	subject.subscribe(client1, key1)

	subject.publish(key2, "msg", false)

	client1.assertReceives("msg", t)
}

func Test11(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("a/**/d")
	key2, _ := parseKey("a/b/c/d/e")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceivesNothing(t)
}

func Test12(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("a/**/d/*")
	key2, _ := parseKey("a/b/c/d/e")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceives("msg", t)
}

func Test13(t *testing.T) {

	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey("a/**/d/f")
	key2, _ := parseKey("a/b/c/d/e")

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	client1.assertReceivesNothing(t)
}

func mockClient(id int) *client {
	client := &client{
		id:    id,
		sendc: make(chan (string), 32),
	}
	return client
}

func (client *client) assertReceives(msg string, t *testing.T) {
	select {
	case actual := <-client.sendc:
		if actual != msg {
			t.Errorf("#%v expected %v got %v", client.id, msg, actual)
		}
		break
	case <-time.After(500):
		t.Errorf("#%v expected %v got nothing", client.id, msg)
		break
	}
}

func (client *client) assertReceivesNothing(t *testing.T) {
	select {
	case actual := <-client.sendc:
		t.Errorf("#%v expected nothing got %v", client.id, actual)
		break
	case <-time.After(500):
		break
	}
}
