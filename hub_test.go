package bus

import (
	"testing"
)

func TestClientRemoved1(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A")
}

func TestClientRemoved2(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A", "B")
}

func TestClientRemoved3(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A", "A/B")
}

func TestClientRemoved4(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "*")
}

func TestClientRemoved5(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "**")
}

func TestClientRemoved6(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "*", "*/A")
}

func TestClientRemoved7(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A", "*/A")
}

func TestClientRemoved8(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A", "A/*")
}

func TestClientRemoved9(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "**", "**/A")
}

func TestClientRemoved10(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A", "**/A")
}

func TestClientRemoved11(t *testing.T) {
	CheckAllSubscriptionsRemoved(t, "A", "A/**")
}

func CheckAllSubscriptionsRemoved(t *testing.T, keys ...string) {
	var subject = newHub()
	var client = mockClient(1)

	for _, str := range keys {
		key, _ := parseKey(str)
		subject.subscribe(client, key)
	}

	subject.deleteSubscriber(client)

	assertHasNoSubscriptions(subject.keyTree, client, t)
}

func assertHasNoSubscriptions(keyTree *keyTree, client *clientHandler, t *testing.T) {
	_, found := keyTree.subscribers[client]
	if found {
		t.Errorf("Expected no subscriptions for client #%v but found one", client.id)
		return
	}

	for _, child := range keyTree.children {
		assertHasNoSubscriptions(child, client, t)
	}
}
