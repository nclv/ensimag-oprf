package main

import (
	"fmt"
	"log"
)

type Key interface {
	Serialize() ([]byte, error)
}

func printKey(key Key) {
	bytesKey, err := key.Serialize()
	if err != nil {
		log.Println(err)
	}

	printByteArray(bytesKey)
}

func printByteArray(byteArray []byte)  {
	log.Println(fmt.Sprintf("%x", byteArray))
}