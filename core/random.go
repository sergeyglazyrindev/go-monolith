package core

import "math/rand"

var OnlyLetersNumbersStringAlphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRandomString(n int, alphabets ...*[]byte) string {
	alphabet := letterRunes
	if len(alphabets) > 0 {
		alphabet = *alphabets[0]
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
