package bll

import (
	"fmt"
	"math"
	"strings"
)

func EnforceHttpPath(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// for our purpose we use specific alphabet to discard some restricted in URL symbols
const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func EncodeBase62(number uint64) string {
	length := len(alphabet)
	var encodedString strings.Builder
	encodedString.Grow(10)
	for ; number > 0; number /= uint64(length) {
		encodedString.WriteByte(alphabet[number%uint64(length)])
	}

	return encodedString.String()
}

func DecodeBase62(encodedString string) (uint64, error) {
	var number uint64
	length := len(alphabet)

	for i, sym := range encodedString {
		alphPos := strings.IndexRune(alphabet, sym)
		if alphPos == -1 {
			return uint64(alphPos), fmt.Errorf("couldn't find symbol:%c in alphabet", sym)
		}
		number += uint64(alphPos) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
