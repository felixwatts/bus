package bus

type keyTree struct {
	children    map[string]*keyTree
	subscribers map[*clientHandler]bool
}

func newKeyTree() *keyTree {
	return &keyTree{
		children:    make(map[string]*keyTree),
		subscribers: make(map[*clientHandler]bool),
	}
}

func overlaps(k1 key, k2 key, k1Doublewild bool, k2Doublewild bool) bool {
	l1 := len(k1)
	l2 := len(k2)

	if l1 == 0 || l2 == 0 {
		return true
	}

	p1 := k1[0]
	p2 := k2[0]

	k1dw := p1 == KEY_DOUBLE_WILD
	k2dw := p2 == KEY_DOUBLE_WILD

	if p1 == KEY_DOUBLE_WILD {
		return overlaps(k1[1:], k2[1:], true, k2dw)
	}

	if k1dw || k2dw || p1 == KEY_WILDCARD || p2 == KEY_WILDCARD || p1 == p2 {
		return overlaps(k1[1:], k2[1:], k1dw || k2Doublewild, k2dw || k2Doublewild)
	}

	return false
}

func (keyTree *keyTree) subscribe(subscriber *clientHandler, key key) *keyTree {
	if len(key) == 0 {
		_, alreadySubscribed := keyTree.subscribers[subscriber]

		keyTree.subscribers[subscriber] = true

		if !alreadySubscribed {
			return keyTree
		}

		return nil
	}

	child, found := keyTree.children[key[0]]
	if !found {
		child = newKeyTree()
		keyTree.children[key[0]] = child
	}

	return child.subscribe(subscriber, key[1:])
}

func (keyTree *keyTree) unsubscribe(subscriber *clientHandler, key key) *keyTree {
	if len(key) == 0 {
		_, isSubscribed := keyTree.subscribers[subscriber]

		if isSubscribed {
			delete(keyTree.subscribers, subscriber)
			return keyTree
		}

		return nil
	}

	child, found := keyTree.children[key[0]]
	if found {
		return child.unsubscribe(subscriber, key[1:])
	}

	return nil
}

func (keyTree *keyTree) publish(key key, msg string, doubleWild bool) {

	if len(key) == 0 {
		for subscriber, _ := range keyTree.subscribers {
			subscriber.send(msg)
		}

		return
	}

	if key[0] == KEY_DOUBLE_WILD {
		keyTree.publish(key[1:], msg, doubleWild)

		for _, child := range keyTree.children {

			child.publish(key, msg, doubleWild)
		}

		return
	}

	if doubleWild {

		keyTree.publish(key[1:], msg, true)

		for childKey, child := range keyTree.children {

			if childKey == key[0] {
				child.publish(key[1:], msg, false)
			}
		}

		return
	}

	if key[0] == KEY_WILDCARD {

		for _, child := range keyTree.children {
			child.publish(key[1:], msg, doubleWild)
		}

		return
	}

	child, found := keyTree.children[key[0]]
	if found {
		child.publish(key[1:], msg, doubleWild)
	}

	child, found = keyTree.children[KEY_WILDCARD]
	if found {
		child.publish(key[1:], msg, doubleWild)
	}

	child, found = keyTree.children[KEY_DOUBLE_WILD]
	if found {
		child.publish(key[1:], msg, true)
	}
}
