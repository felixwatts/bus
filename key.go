package bus

import (
	"errors"
	"strings"
)

type key []string

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
