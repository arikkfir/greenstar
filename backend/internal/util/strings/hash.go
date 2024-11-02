package strings

import (
	"math/rand"
)

var (
	hashLetters = []rune("abcdefghijklmnopqrstuvwxyz")
)

func RandomHash(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = hashLetters[rand.Intn(len(hashLetters))]
	}
	return string(b)
}
