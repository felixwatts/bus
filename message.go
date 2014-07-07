package bus

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MSG_PART_SEPARATOR = " "
	KEY_PART_SEPARATOR = "/"
	KEY_WILDCARD       = "*"
	KEY_DOUBLE_WILD    = "**"
	MSG_TYPE_OK        = "OK"
	MSG_TYPE_SUBSCRIBE = "SB"
	MSG_TYPE_PUBLISH   = "PB"
	MSG_TYPE_FAIL      = "FL"
)

type message struct {
	meaning   string
	requestId string
	key       key
	val       string
}

type key []string

func (message message) String() string {
	if message.val != "" {
		return fmt.Sprintf("%v%v%v%v%v%v%v\n", message.meaning, MSG_PART_SEPARATOR, message.requestId, MSG_PART_SEPARATOR, message.key, MSG_PART_SEPARATOR, message.val)
	} else if message.key != nil {
		return fmt.Sprintf("%v%v%v%v%v\n", message.meaning, MSG_PART_SEPARATOR, message.requestId, MSG_PART_SEPARATOR, message.key)
	} else if message.requestId != "" {
		return fmt.Sprintf("%v%v%v\n", message.meaning, MSG_PART_SEPARATOR, message.requestId)
	} else {
		return fmt.Sprintf("%v\n", message.meaning)
	}
}

func (key key) String() string {
	return strings.Join(key, KEY_PART_SEPARATOR)
}

func parseKey(keyStr string) (key, error) {
	keyParts := key(strings.Split(keyStr, KEY_PART_SEPARATOR))

	if len(keyParts) == 0 {
		return nil, errors.New("Empty key")
	}

	for _, keyPart := range keyParts {

		if keyPart == KEY_WILDCARD || keyPart == KEY_DOUBLE_WILD {
			continue
		}

		if strings.Contains(keyPart, KEY_WILDCARD) {
			return nil, errors.New("Invalid key: Part contains wildcard.")
		}

		if len(keyPart) == 0 {
			return nil, errors.New("Invalid key: Empty part")
		}
	}

	return keyParts, nil
}

func parseMessage(input string) (message, error) {
	parts := strings.Split(input, MSG_PART_SEPARATOR)
	if len(parts) == 0 {
		return message{}, errors.New("Empty message.")
	}

	if len(parts) > 4 {
		return message{}, errors.New("Invalid message: Too many parts.")
	}

	result := message{
		meaning: parts[0],
	}

	if len(parts) > 1 {
		result.requestId = parts[1]
	}

	if len(parts) > 2 {
		key, err := parseKey(parts[2])
		if err != nil {
			return message{}, err
		}
		result.key = key
	}

	if len(parts) > 3 {
		result.val = parts[3]
	}

	return result, nil
}
