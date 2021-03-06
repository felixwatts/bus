package bus

type hub struct {
	subscriptionBySubscriber map[*clientHandler]map[*keyTree]bool
	keyTree                  *keyTree
}

func newHub() *hub {
	return &hub{
		subscriptionBySubscriber: make(map[*clientHandler]map[*keyTree]bool),
		keyTree:                  newKeyTree(),
	}
}

func (hub *hub) getClientSubscriptions(subscriber *clientHandler) map[*keyTree]bool {
	result, found := hub.subscriptionBySubscriber[subscriber]
	if !found {
		result = make(map[*keyTree]bool)
		hub.subscriptionBySubscriber[subscriber] = result
	}
	return result
}

func (hub *hub) subscribe(subscriber *clientHandler, key key) {
	var sub = hub.keyTree.subscribe(subscriber, key)

	if sub != nil {
		hub.getClientSubscriptions(subscriber)[sub] = true
	}
}

func (hub *hub) unsubscribe(subscriber *clientHandler, key key) {
	var sub = hub.keyTree.unsubscribe(subscriber, key)

	if sub != nil {
		delete(hub.getClientSubscriptions(subscriber), sub)
	}
}

func (hub *hub) publish(publisher *clientHandler, key key, msg string) bool {

	subscribers := make(map[*clientHandler]bool)
	canPublish := hub.keyTree.findSubscribers(key, false, publisher, subscribers)

	if !canPublish {
		return false
	}

	for subscriber, _ := range subscribers {
		subscriber.send(msg)
	}

	return true
}

func (hub *hub) deleteSubscriber(subscriber *clientHandler) {
	var subscriptions = hub.getClientSubscriptions(subscriber)

	for sub, _ := range subscriptions {
		delete(sub.subscribers, subscriber)
	}
}

func (hub *hub) claim(claimant *clientHandler, key key) bool {
	return hub.keyTree.claim(claimant, key)
}
