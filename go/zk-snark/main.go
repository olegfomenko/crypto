// Package zk_snark
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package zk_snark

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/cloudflare/bn256"
)

type (
	// Public functions

	L1 func(xi []*bn256.G1) []*bn256.G1
	L2 func(xi []*bn256.G2) []*bn256.G2
	R2 func(xi []*bn256.G2) []*bn256.G2
	O2 func(xi []*bn256.G2) []*bn256.G2

	// Private functions

	BigL1 func(xi []*bn256.G1) *bn256.G1
	BigL2 func(xi []*bn256.G2) *bn256.G2
	BigR2 func(xi []*bn256.G2) *bn256.G2
	BigO2 func(xi []*bn256.G2) *bn256.G2
	H2    func(xi []*bn256.G2) *bn256.G2

	SetupParams struct {
		G1 *bn256.G1
		G2 *bn256.G2

		G1_ts    *bn256.G1
		G1_alpha *bn256.G1
		G1_si    []*bn256.G1
		G2_si    []*bn256.G2

		G1_l []*bn256.G1
		G2_l []*bn256.G2
		G2_r []*bn256.G2
		G2_o []*bn256.G2

		G2_alpha_l []*bn256.G2
		G2_alpha_r []*bn256.G2
		G2_alpha_o []*bn256.G2

		N uint64
	}

	Proof struct {
		G1_L       *bn256.G1
		G2_L       *bn256.G2
		G2_R       *bn256.G2
		G2_O       *bn256.G2
		G2_alpha_L *bn256.G2
		G2_alpha_R *bn256.G2
		G2_alpha_O *bn256.G2
		G2_h       *bn256.G2
	}
)

func Setup(l1 L1, l2 L2, r R2, o O2, n uint64) *SetupParams {
	s, alpha := GetRandomWithMax(bn256.Order), GetRandomWithMax(bn256.Order)
	defer func() {
		s.SetUint64(0)
		alpha.SetUint64(0)
	}()

	_, g1, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		panic(err)
	}

	_, g2, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		panic(err)
	}

	g1_si := make([]*bn256.G1, 0, n)
	g2_si := make([]*bn256.G2, 0, n)
	g2_alphasi := make([]*bn256.G2, 0, n)

	for i := uint64(0); i < n; i++ {
		si := new(big.Int).Exp(s, new(big.Int).SetUint64(i), bn256.Order)

		g1_si = append(g1_si, new(bn256.G1).ScalarMult(g1, si))
		g2_si = append(g2_si, new(bn256.G2).ScalarMult(g2, si))

		alphasi := mul(alpha, si)
		g2_alphasi = append(g2_alphasi, new(bn256.G2).ScalarMult(g2, alphasi))
	}

	return &SetupParams{
		G1:       g1,
		G2:       g2,
		G1_ts:    new(bn256.G1).ScalarMult(g1, t(s, n)),
		G1_alpha: new(bn256.G1).ScalarMult(g1, alpha),
		G1_si:    g1_si,
		G2_si:    g2_si,

		G1_l: l1(g1_si),
		G2_l: l2(g2_si),
		G2_r: r(g2_si),
		G2_o: o(g2_si),

		G2_alpha_l: l2(g2_alphasi),

		G2_alpha_r: r(g2_alphasi),
		G2_alpha_o: o(g2_alphasi),
		N:          n,
	}
}

func MakeProof(params *SetupParams, bigL1 BigL1, bigL2 BigL2, bigR BigR2, bigO BigO2, h H2) *Proof {
	return &Proof{
		G1_L:       bigL1(params.G1_l),
		G2_L:       bigL2(params.G2_l),
		G2_alpha_L: bigL2(params.G2_alpha_l),
		G2_R:       bigR(params.G2_r),
		G2_alpha_R: bigR(params.G2_alpha_r),
		G2_O:       bigO(params.G2_o),
		G2_alpha_O: bigO(params.G2_alpha_o),
		G2_h:       h(params.G2_si),
	}
}

func VerifyProof(params *SetupParams, proof *Proof) error {
	e1 := bn256.Pair(params.G1, proof.G2_alpha_L)
	e2 := bn256.Pair(params.G1_alpha, proof.G2_L)
	if !bytes.Equal(e1.Marshal(), e2.Marshal()) {
		return errors.New("check #1 failed")
	}

	e1 = bn256.Pair(params.G1, proof.G2_alpha_R)
	e2 = bn256.Pair(params.G1_alpha, proof.G2_R)
	if !bytes.Equal(e1.Marshal(), e2.Marshal()) {
		return errors.New("check #2 failed")
	}

	e1 = bn256.Pair(params.G1, proof.G2_alpha_O)
	e2 = bn256.Pair(params.G1_alpha, proof.G2_O)
	if !bytes.Equal(e1.Marshal(), e2.Marshal()) {
		return errors.New("check #3 failed")
	}

	e1 = bn256.Pair(proof.G1_L, proof.G2_R)
	e2 = bn256.Pair(params.G1_ts, proof.G2_h)
	e3 := bn256.Pair(params.G1, proof.G2_O)

	if !bytes.Equal(e1.Marshal(), new(bn256.GT).Add(e2, e3).Marshal()) {
		return errors.New("check #4 failed")
	}

	return nil
}

func t(x *big.Int, n uint64) *big.Int {
	res := new(big.Int).SetUint64(1)

	for i := uint64(1); i <= n; i++ {
		res = mul(res, sub(x, new(big.Int).SetUint64(i)))
	}

	return res
}

func mul(a, b *big.Int) *big.Int {
	return mod(new(big.Int).Mul(a, b))
}

func sub(a, b *big.Int) *big.Int {
	return mod(new(big.Int).Sub(a, b))
}

func add(a, b *big.Int) *big.Int {
	return mod(new(big.Int).Sub(a, b))
}

func mod(v *big.Int) *big.Int {
	return new(big.Int).Mod(v, bn256.Order)
}
