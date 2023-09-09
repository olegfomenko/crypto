// Package pedersen
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type SchnorrSignature struct {
	R *bn256.G1
	S *big.Int
}

func SignSchnorr(prv *big.Int, publicKey *bn256.G1, m *big.Int) (SchnorrSignature, error) {
	k, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		return SchnorrSignature{}, err
	}

	kG := ScalarMul(G, k)
	hash := Hash(m.Bytes(), X(publicKey).Bytes(), Y(publicKey).Bytes())
	s := add(k, minus(mul(hash, prv)))

	return SchnorrSignature{
		R: kG,
		S: s,
	}, nil
}

func VerifySchnorr(sig SchnorrSignature, publicKey *bn256.G1, m *big.Int) error {
	hash := Hash(m.Bytes(), X(publicKey).Bytes(), Y(publicKey).Bytes())

	p1 := ScalarMul(publicKey, hash)
	p1 = Sub(sig.R, p1)

	p2 := ScalarMul(G, sig.S)

	if !bytes.Equal(p1.Marshal(), p2.Marshal()) {
		return errors.New("verification failed")
	}

	return nil
}
