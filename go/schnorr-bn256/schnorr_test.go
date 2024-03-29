// Package schnorr_bn256
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package schnorr_bn256

import (
	"crypto/rand"
	"fmt"
	"github.com/cloudflare/bn256"
	"math/big"
	"testing"
)

func TestSchnorrSignature(t *testing.T) {
	_, G, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		panic(err)
	}

	alicePrv, alicePub, err := R(G)
	if err != nil {
		panic(err)
	}

	bobPrv, bobPub, err := R(G)
	if err != nil {
		panic(err)
	}

	alice_r, aliceR, err := R(G)
	if err != nil {
		panic(err)
	}

	bob_r, bobR, err := R(G)
	if err != nil {
		panic(err)
	}

	R := new(bn256.G1).Add(aliceR, bobR)

	Pub := new(bn256.G1).Add(alicePub, bobPub)

	message := Msg([]byte("Hello world"))
	fmt.Println(message.String())

	aliceSig, err := MultiSigSchnorr(alicePrv, alice_r, Pub, R, message)
	if err != nil {
		panic(err)
	}

	bobSig, err := MultiSigSchnorr(bobPrv, bob_r, Pub, R, message)
	if err != nil {
		panic(err)
	}

	sig := &SchnorrSignature{
		R: aliceSig.R,
		S: new(big.Int).Mod(new(big.Int).Add(aliceSig.S, bobSig.S), bn256.Order),
	}

	if !VerifySchnorr(sig, Pub, G, message) {
		panic("failed")
	}
}
