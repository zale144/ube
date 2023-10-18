package model

type Entity interface {
	GetKey() Key
}

type Key interface {
	PK() string
}

type SK interface {
	SK() string
}

const KeySeparator = "|"

func StringifyKey(key Key) string {
	keyStr := key.PK()
	if sk, ok := key.(SK); ok {
		keyStr = keyStr + KeySeparator + sk.SK()
	}

	return keyStr
}
