// Package pedersen
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

var G *bn256.G1
var H *bn256.G1

func ScalarMul(p *bn256.G1, k *big.Int) *bn256.G1 {
	return new(bn256.G1).ScalarMult(p, k)
}

func Add(a, b *bn256.G1) *bn256.G1 {
	return new(bn256.G1).Add(a, b)
}

func Sub(a, b *bn256.G1) *bn256.G1 {
	return Add(a, new(bn256.G1).Neg(b))
}

func X(p *bn256.G1) *big.Int {
	bytes := p.Marshal()
	return new(big.Int).SetBytes(bytes[:32])
}

func Y(p *bn256.G1) *big.Int {
	bytes := p.Marshal()
	return new(big.Int).SetBytes(bytes[32:])
}
