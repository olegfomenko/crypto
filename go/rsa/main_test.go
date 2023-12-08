// Package rsa
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package rsa

import (
	"fmt"
	"math/big"
	"testing"
)

func TestRSA(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	msg := new(big.Int).SetBytes([]byte("Hello World"))
	cypher := Encrypt(msg, key.PublicKey)
	text := Decrypt(cypher, key)
	fmt.Println("Original message:", string(msg.Bytes()))
	fmt.Println("Encrypted message:", cypher)
	fmt.Println("Decrypted message:", string(text.Bytes()))
}
