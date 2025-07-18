// Package rsa
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package rsa

import (
	"math/big"

	"github.com/olegfomenko/crypto/go/math"
)

const size = 256

const Exp = 65537

var e = big.NewInt(Exp)

type PublicKey struct {
	n *big.Int
}

type PrivateKey struct {
	*PublicKey
	p, q *big.Int
	phiN *big.Int
	d    *big.Int
}

func GeneratePrivateKey() (*PrivateKey, error) {
	p, err := math.GenRandPrime(size)
	if err != nil {
		return nil, err
	}

	q, err := math.GenRandPrime(size)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).Mul(p, q)

	// phi(n) = (p - 1)(q - 1) because p and q - primes
	phiN := new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1)))

	// Euler's theorem can be used
	d := new(big.Int).ModInverse(e, phiN)

	return &PrivateKey{
		PublicKey: &PublicKey{
			n: n,
		},
		p:    p,
		q:    q,
		phiN: phiN,
		d:    d,
	}, nil
}

func Encrypt(msg *big.Int, pk *PublicKey) *big.Int {
	return new(big.Int).Exp(msg, e, pk.n)
}

func Decrypt(cypher *big.Int, prv *PrivateKey) *big.Int {
	return new(big.Int).Exp(cypher, prv.d, prv.n)
}
