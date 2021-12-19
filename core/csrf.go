package core

import (
	"bytes"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_-+=][{}'\\\"|;:/?.>,<`")
var csrfAlphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRandomStringForCSRFToken() string {
	n := 32
	b := make([]byte, n)
	for i := range b {
		b[i] = csrfAlphabet[rand.Intn(len(csrfAlphabet))]
	}
	return string(b)
}

func GenerateCSRFToken() string {
	return GenerateRandomStringForCSRFToken()
}

func MaskCSRFToken(token string) string {
	mask := GenerateRandomStringForCSRFToken()
	cipher := ""
	for i := range token {
		tokenIndex := bytes.IndexByte(csrfAlphabet, token[i])
		maskIndex := bytes.IndexByte(csrfAlphabet, mask[i])
		cipher = cipher + string(csrfAlphabet[(tokenIndex+maskIndex)%len(csrfAlphabet)])
	}
	return mask + cipher
}

func UnmaskCSRFToken(token string) string {
	tokenpart := token[32:]
	maskpart := token[:32]
	ret := ""
	for i := range tokenpart {
		tokenIndex := bytes.IndexByte(csrfAlphabet, tokenpart[i])
		maskIndex := bytes.IndexByte(csrfAlphabet, maskpart[i])
		if tokenIndex >= maskIndex {
			ret = ret + string(csrfAlphabet[tokenIndex-maskIndex])
		} else {
			ret = ret + string(csrfAlphabet[len(csrfAlphabet)+(tokenIndex-maskIndex)])
		}
	}
	return ret
}
