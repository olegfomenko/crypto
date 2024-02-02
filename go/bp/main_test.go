// Package bp
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bp

import (
	"math/big"
	"testing"
)

func TestBulletProof(t *testing.T) {
	const n = 4
	public := NewBulletProofPublic(n)

	proof, err := public.Prove(big.NewInt(11), values(1)[0])
	if err != nil {
		panic(err)
	}

	err = public.Verify(proof)
	if err != nil {
		panic(err)
	}
}

func TestInnerProduct(t *testing.T) {
	const n = 4
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
