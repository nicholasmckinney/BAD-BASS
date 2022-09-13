package internal

import (
	"crypto/rand"
	"errors"
	"os"
)

const (
	Alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func RandomString(length int) string {
	ll := len(Alphanumeric)
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = Alphanumeric[int(b[i])%ll]
	}
	return string(b)
}

func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
