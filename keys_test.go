package bus

import (
	"testing"
	"time"
)

var subject *keyTree

func Test2(t *testing.T) {
	assertKeysMatch("a", "a", true, t)
}

func Test3(t *testing.T) {
	assertKeysMatch("a", "b", false, t)
}

func Test4(t *testing.T) {
	assertKeysMatch("a/b", "a/b", true, t)
}

func Test5(t *testing.T) {
	assertKeysMatch("a/*", "a/b", true, t)
}

func Test6(t *testing.T) {
	assertKeysMatch("*/b", "a/b", true, t)
}

func Test7(t *testing.T) {
	assertKeysMatch("*/c", "a/b", false, t)
}

func Test8(t *testing.T) {
	assertKeysMatch("c/b", "a/b", false, t)
}

func Test9(t *testing.T) {
	assertKeysMatch("a/**", "a/b/c", true, t)
}

func Test10(t *testing.T) {
	assertKeysMatch("a/**/d", "a/b/c/d", true, t)
}

func Test11(t *testing.T) {
	assertKeysMatch("a/**/d", "a/b/c/d/e", false, t)
}

func Test12(t *testing.T) {
	assertKeysMatch("a/**/d/*", "a/b/c/d/e", true, t)
}

func Test13(t *testing.T) {
	assertKeysMatch("a/**/d/f", "a/b/c/d/e", false, t)
}

func Test14(t *testing.T) {
	assertKeysMatch("a", "*", true, t)
}

func Test15(t *testing.T) {
	assertKeysMatch("a/b", "*", false, t)
}

func Test16(t *testing.T) {
	assertKeysMatch("a/b", "a/*", true, t)
}

func Test17(t *testing.T) {
	assertKeysMatch("a/b", "*/b", true, t)
}

func Test18(t *testing.T) {
	assertKeysMatch("a/b", "*/c", false, t)
}

func Test19(t *testing.T) {
	assertKeysMatch("a/b", "*/*", true, t)
}

func assertKeysMatch(subKey string, pubKey string, match bool, t *testing.T) {
	subject = newKeyTree()

	client1 := mockClient(1)

	key1, _ := parseKey(subKey)
	key2, _ := parseKey(pubKey)

	subject.subscribe(client1, key1)
	subject.publish(key2, "msg", false)

	if match {
		client1.assertReceives("msg", t)
	} else {
		client1.assertReceivesNothing(t)
	}
}

func mockClient(id int) *clientHandler {
	client := &clientHandler{
		id:    id,
		sendc: make(chan (string), 32),
	}
	return client
}

func (client *clientHandler) assertReceives(msg string, t *testing.T) {
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

func (client *clientHandler) assertReceivesNothing(t *testing.T) {
	select {
	case actual := <-client.sendc:
		t.Errorf("#%v expected nothing got %v", client.id, actual)
		break
	case <-time.After(500):
		break
	}
}
