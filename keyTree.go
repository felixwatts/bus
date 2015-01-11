package bus

type keyTree struct {
	owner       *clientHandler
	children    map[string]*keyTree
	subscribers map[*clientHandler]bool
}

func newKeyTree() *keyTree {
	return &keyTree{
		children:    make(map[string]*keyTree),
		subscribers: make(map[*clientHandler]bool),
	}
}

func (keyTree *keyTree) claim(claimant *clientHandler, key key) bool {

	if keyTree.owner != nil && keyTree.owner != claimant {
		return false
	}

	if len(key) == 0 {
		keyTree.owner = claimant
		return true
	}

	if key[0] == KEY_WILDCARD || key[0] == KEY_DOUBLE_WILD {
		return false
	}

	child, found := keyTree.children[key[0]]
	if !found {
		child = newKeyTree()
		keyTree.children[key[0]] = child
	}

	return child.claim(claimant, key[1:])
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

func (keyTree *keyTree) findSubscribers(key key, doubleWild bool, publisher *clientHandler, subscribers map[*clientHandler]bool) bool {
	if keyTree.owner != nil && keyTree.owner != publisher {
		return false
	}

	if len(key) == 0 {
		for subscriber, _ := range keyTree.subscribers {
			subscribers[subscriber] = true
		}

		return true
	}

	isAllowed := true

	if key[0] == KEY_DOUBLE_WILD {
		isAllowed = keyTree.findSubscribers(key[1:], doubleWild, publisher, subscribers)

		if !isAllowed {
			return false
		}

		for _, child := range keyTree.children {

			isAllowed = child.findSubscribers(key, doubleWild, publisher, subscribers)

			if !isAllowed {
				return false
			}
		}

		return true
	}

	if doubleWild {

		isAllowed = keyTree.findSubscribers(key[1:], true, publisher, subscribers)

		if !isAllowed {
			return false
		}

		for childKey, child := range keyTree.children {

			if childKey == key[0] {
				isAllowed = child.findSubscribers(key[1:], false, publisher, subscribers)

				if !isAllowed {
					return false
				}
			}
		}

		return true
	}

	if key[0] == KEY_WILDCARD {

		for _, child := range keyTree.children {
			isAllowed = child.findSubscribers(key[1:], doubleWild, publisher, subscribers)

			if !isAllowed {
				return false
			}
		}

		return true
	}

	child, found := keyTree.children[key[0]]
	if found {
		isAllowed = child.findSubscribers(key[1:], doubleWild, publisher, subscribers)

		if !isAllowed {
			return false
		}
	}

	child, found = keyTree.children[KEY_WILDCARD]
	if found {
		isAllowed = child.findSubscribers(key[1:], doubleWild, publisher, subscribers)

		if !isAllowed {
			return false
		}
	}

	child, found = keyTree.children[KEY_DOUBLE_WILD]
	if found {
		isAllowed = child.findSubscribers(key[1:], true, publisher, subscribers)

		if !isAllowed {
			return false
		}
	}

	return true
}
