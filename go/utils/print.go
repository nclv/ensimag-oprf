package utils

import (
	"fmt"
	"log"
)

type Key interface {
	Serialize() ([]byte, error)
}

func PrintKey(key Key) {
	bytesKey, err := key.Serialize()
	if err != nil {
		log.Println(err)
	}

	PrintByteArray(bytesKey)
}

func PrintByteArray(byteArray []byte) {
	log.Println(fmt.Sprintf("%x", byteArray))
}
