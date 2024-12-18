package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math"
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

func get_first_k_bits(k int, hash []byte) uint64 {
	bs := NewBStreamReader(hash)
	k_bits, err := bs.ReadBits(k)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		panic(1)
	}

	return k_bits
}

func get_msb_of_remaining_hash(k int, hash []byte) int {
	msb := 0
	bs := NewBStreamReader(hash)
	for msb <= len(hash)*8 {
		bt, err := bs.ReadBit()
		if err != nil {
			fmt.Println("Error: ", err)
			panic(1)
		}
		msb++
		// fmt.Println(bt)
		if bt && msb > k {
			return msb - k
		}
	}
	return msb - k
}

func hyperLogLog(k int, elems []string) uint64 {
	cardinality := uint64(0)
	buckets_len := int(math.Pow(2, float64(k)))
	buckets := make([]uint64, buckets_len)

	for _, v := range elems {
		hash := generate_hash(v)
		bucket_idx := get_first_k_bits(k, hash)
		msb := get_msb_of_remaining_hash(k, hash)
		buckets[bucket_idx] = max(buckets[bucket_idx], uint64(msb))
	}

	Z := float64(0)
	for _, v := range buckets {
		Z += math.Pow(float64(1/2), float64(v))
	}

	cardinality = uint64(float64(0.79402) * float64(buckets_len) * float64(buckets_len) * Z)

	return cardinality
}

func main() {
	// Runner function

	// This example outputs 203
	// cardinality := hyperLogLog(3, []string{"hi", "bye", "yo", "hi", "hello", "main", "next"})
	// fmt.Println(cardinality)

	return
}
