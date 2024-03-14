// Package bp
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bp

import (
	"crypto/rand"
	"github.com/davecgh/go-spew/spew"
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

func TestBulletProof(t *testing.T) {
	const n = 64
	public := NewBulletProofPublic(n)

	prv, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}

	proof, err := public.Prove(big.NewInt(11), prv)
	if err != nil {
		panic(err)
	}

	spew.Dump(proof)

	err = public.Verify(proof)
	if err != nil {
		panic(err)
	}
}

func TestInnerProduct(t *testing.T) {
	const n = 8
	public := NewInnerArgumentPublic(n)

	a := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(10), big.NewInt(1)}
	b := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(10)}

	proof, err := public.Prove(a, b)
	if err != nil {
		panic(err)
	}

	err = public.Verify(proof)
	if err != nil {
		panic(err)
	}
}
