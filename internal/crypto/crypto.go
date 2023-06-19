package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func Hash(message string, key string) string {
	if key == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	dst := mac.Sum(nil)

	return fmt.Sprintf("%x", dst)
}

func ValidMAC(message, messageMAC, key string) bool {
	if key == "" && messageMAC == "" {
		return true
	}
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := fmt.Sprintf("%x", mac.Sum(nil))
	if messageMAC == expectedMAC {
		return true
	} else {
		return false
	}
}
