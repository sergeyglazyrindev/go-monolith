package services

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
