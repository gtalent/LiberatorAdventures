package main

import (
	"crypto/sha512"
)

func PasswordHash(password string) []byte {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hasher.Sum()
}
