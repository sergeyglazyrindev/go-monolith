package services

import (
	"crypto/rand"
	"math/big"
)

// GenerateBase64 generates a base64 string of length length
func GenerateBase64(length int) string {
	base := new(big.Int)
	base.SetString("64", 10)

	base64 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	tempKey := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, base)
		tempKey += string(base64[int(index.Int64())])
	}
	return tempKey
}

//// GenerateBase32 generates a base64 string of length length
//func GenerateBase32(length int) string {
//	base := new(big.Int)
//	base.SetString("32", 10)
//
//	base32 := "234567abcdefghijklmnopqrstuvwxyz"
//	tempKey := ""
//	for i := 0; i < length; i++ {
//		index, _ := rand.Int(rand.Reader, base)
//		tempKey += string(base32[int(index.Int64())])
//	}
//	return tempKey
//}
//
