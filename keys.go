package bus

type keyTree struct {
	children    map[string]*keyTree
	subscribers []*client
}

func newKeyTree() *keyTree {
	return &keyTree{
		children:    make(map[string]*keyTree),
		subscribers: make([]*client, 0),
	}
}

func (keyTree *keyTree) subscribe(subscriber *client, key key) {
	if len(key) == 0 {
		keyTree.subscribers = append(keyTree.subscribers, subscriber)
		return
	}

	child, found := keyTree.children[key[0]]
	if !found {
		child = newKeyTree()
		keyTree.children[key[0]] = child
	}

	child.subscribe(subscriber, key[1:])
}

func (keyTree *keyTree) publish(key key, msg string, doubleWild bool) {

	if len(key) == 0 {
		for _, subscriber := range keyTree.subscribers {
			subscriber.send(msg)
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
