package utils

import (
	"crypto/rand"
	"io"
)

const codeLenght = 6

var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateCode() string {
	b := make([]byte, codeLenght)
	io.ReadAtLeast(rand.Reader, b, codeLenght)

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
