package main

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// The bitstream implementation taken from https://github.com/kkdai/bstream
type bit bool

const (
	zero bit = false
	one  bit = true
)

type BStream struct {
	stream []byte
	rCount uint8
}

func NewBStreamReader(data []byte) *BStream {
	return &BStream{stream: data, rCount: 8}
}

func (b *BStream) ReadBit() (bit, error) {
	if len(b.stream) == 0 {
		return zero, io.EOF
	}

	if b.rCount == 0 {
		b.stream = b.stream[1:]

		if len(b.stream) == 0 {
			return zero, io.EOF
		}

		b.rCount = 8
	}

	retBit := b.stream[0] & (1 << (b.rCount - 1))
	b.rCount--

	return retBit != 0, nil
}

func (b *BStream) ReadByte() (byte, error) {
	if len(b.stream) == 0 {
		return 0, io.EOF
	}

	if b.rCount == 0 {
		b.stream = b.stream[1:]

		if len(b.stream) == 0 {
			return 0, io.EOF
		}

		b.rCount = 8
	}

	if b.rCount == 8 {
		byt := b.stream[0]
		b.stream = b.stream[1:]
		return byt, nil
	}

	retByte := b.stream[0] << (8 - b.rCount)
	b.stream = b.stream[1:]

	if len(b.stream) == 0 {
		return 0, io.EOF
	}

	retByte |= b.stream[0] >> b.rCount
	return retByte, nil
}

func (b *BStream) ReadBits(count int) (uint64, error) {
	var retValue uint64

	for count >= 8 {
		retValue <<= 8
		byt, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		retValue |= uint64(byt)
		count = count - 8
	}

	for count > 0 {
		retValue <<= 1
		bi, err := b.ReadBit()
		if err != nil {
			return 0, err
		}
		if bi {
			retValue |= 1
		}

		count--
	}

	return retValue, nil
}

func generate_hash(text string) []byte {
	sha := sha1.New()
	sha.Write([]byte(text))
	hash := sha.Sum(nil)
	return hash
}

func get_first_k_bits(k int, hash []byte) {
	bs := NewBStreamReader(hash)
	k_bits, err := bs.ReadBits(k)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%b\n", k_bits)
}

func main() {
	temp_str := "Hello"
	hash := generate_hash(temp_str)
	fmt.Printf("%b\n", hash)
	get_first_k_bits(10, hash)

	return
}
