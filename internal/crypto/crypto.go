package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Hash(message string, key string) string {
	if key == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	dst := mac.Sum(nil)

	return string(dst)
}

func ValidMAC(message, messageMAC, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal([]byte(messageMAC), expectedMAC)
}
