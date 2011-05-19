/*
 * Copyright 2011 <gtalent2@gmail.com>
 * This file is released under the BSD license, as defined here:
 * 	http://www.opensource.org/licenses/bsd-license.php
 */
package util

import (
	"crypto/sha512"
)

type Password struct {
	Version uint
	Hash    []byte
}

//Returns the password hash and the version number of the hashing strategy.
//Pass in -1 for the current method of password hashing.
func PasswordHash(password string, version int) Password {
	switch version {
	default:
		bytes := []byte(password)
		hasher := sha512.New()
		for i := 0; i < 100; i++ {
			hasher.Write(bytes)
			bytes = hasher.Sum()
		}
		var retval Password
		retval.Version = 0
		retval.Hash = bytes
		return retval
	}
	var retval Password
	return retval
}
