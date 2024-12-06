package main

import (
	"crypto/sha1"
	"fmt"
)

func generate_hash(text string) []byte {
	sha := sha1.New()
	sha.Write([]byte(text))
	hash := sha.Sum(nil)
	return hash
}

func main() {
	fmt.Println("Hello, World!")
	return
}
