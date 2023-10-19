package utils

import (
	"encoding/base64"
)

func EncodeBytesToBase64(bytes []byte) []byte {
	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(encodedBytes, bytes)
	return encodedBytes
}

func DecodeBase64ToBytes(encodedBytes string) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(encodedBytes)
	return bytes, err
}

func EncodeStringToBase64(str string) string {
	return string(EncodeBytesToBase64([]byte(str)))
}

func DecodeBase64ToString(encodedStr string) (string, error) {
	bytes, err := DecodeBase64ToBytes(encodedStr)
	return string(bytes), err
}
