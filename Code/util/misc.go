/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package util

import (
	"crypto/sha512"
)

func PasswordHash(password string) []byte {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hasher.Sum()
}
