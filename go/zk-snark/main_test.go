// Package zk_snark
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package zk_snark

import (
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

// TEST FUNCTION
// f(w,a,b) = w? (a * b) : (a + b)
// v = w(a*b) + (1-w) * (a+b)
// m = a*b
// v = w*m + (1-w) * (a+b)

// Gates
// a b w m v

// Proof that we know `a`, such that f(1, a, 2) = 8
// a = 4

func NewInt(val int64) *big.Int {
	return mod(big.NewInt(val))
}

var inverse2 = new(big.Int).ModInverse(big.NewInt(2), bn256.Order)
var inverse4 = new(big.Int).ModInverse(big.NewInt(4), bn256.Order)

func f1(xi []*bn256.G1, c []*big.Int, inverse *big.Int) *bn256.G1 {
	var e = make([]*bn256.G1, 3)
	for i, val := range xi {
		e[i] = new(bn256.G1).ScalarMult(val, c[i])
	}

	var res = e[0]
	for i := 1; i < 3; i++ {
		res = new(bn256.G1).Add(e[i], res)
	}

	if inverse == nil {
		return res
	}

	return new(bn256.G1).ScalarMult(res, inverse)
}

func f2(xi []*bn256.G2, c []*big.Int, inverse *big.Int) *bn256.G2 {
	var e = make([]*bn256.G2, 0, 3)
	for i, val := range xi {
		e = append(e, new(bn256.G2).ScalarMult(val, c[i]))
	}

	var res = e[0]
	for i := 1; i < 3; i++ {
		res = new(bn256.G2).Add(e[i], res)
	}

	if inverse == nil {
		return res
	}

	return new(bn256.G2).ScalarMult(res, inverse)
}

func l1(xi []*bn256.G1) []*bn256.G1 {
	la := f1(xi, []*big.Int{NewInt(6), NewInt(-5), NewInt(1)}, inverse2)
	lw := f1(xi, []*big.Int{NewInt(-4), NewInt(5), NewInt(-1)}, inverse2)
	return []*bn256.G1{la, nil, lw, nil, nil}
}

func l2(xi []*bn256.G2) []*bn256.G2 {
	la := f2(xi, []*big.Int{NewInt(6), NewInt(-5), NewInt(1)}, inverse2)
	lw := f2(xi, []*big.Int{NewInt(-4), NewInt(5), NewInt(-1)}, inverse2)
	return []*bn256.G2{la, nil, lw, nil, nil}
}

func r2(xi []*bn256.G2) []*bn256.G2 {
	ra := f2(xi, []*big.Int{NewInt(3), NewInt(-4), NewInt(1)}, nil)
	rb := f2(xi, []*big.Int{NewInt(12), NewInt(-13), NewInt(3)}, inverse2)
	rw := f2(xi, []*big.Int{NewInt(2), NewInt(-3), NewInt(1)}, inverse2)
	rm := f2(xi, []*big.Int{NewInt(-3), NewInt(4), NewInt(-1)}, nil)
	return []*bn256.G2{ra, rb, rw, rm, nil}
}

func o2(xi []*bn256.G2) []*bn256.G2 {
	oa := f2(xi, []*big.Int{NewInt(3), NewInt(-4), NewInt(1)}, nil)
	ob := f2(xi, []*big.Int{NewInt(3), NewInt(-4), NewInt(1)}, nil)
	ow := f2(xi, []*big.Int{NewInt(2), NewInt(-3), NewInt(1)}, inverse2)
	om := f2(xi, []*big.Int{NewInt(6), NewInt(-5), NewInt(1)}, inverse2)
	ov := f2(xi, []*big.Int{NewInt(-3), NewInt(4), NewInt(-1)}, nil)
	return []*bn256.G2{oa, ob, ow, om, ov}
}

// a b w m v
var inputGates = []*big.Int{NewInt(4), NewInt(2), NewInt(1), NewInt(8), NewInt(8)}

func big1(xi []*bn256.G1) *bn256.G1 {
	var e = make([]*bn256.G1, 0, 5)
	for i, val := range xi {
		if val != nil {
			e = append(e, new(bn256.G1).ScalarMult(val, inputGates[i]))
		}
	}

	var res = e[0]
	for i := 1; i < len(e); i++ {
		res = new(bn256.G1).Add(e[i], res)
	}

	return res
}

func big2(xi []*bn256.G2) *bn256.G2 {
	var e = make([]*bn256.G2, 0, 5)
	for i, val := range xi {
		if val != nil {
			e = append(e, new(bn256.G2).ScalarMult(val, inputGates[i]))
		}
	}

	var res = e[0]
	for i := 1; i < len(e); i++ {
		res = new(bn256.G2).Add(e[i], res)
	}

	return res
}

func h(xi []*bn256.G2) *bn256.G2 {
	c := []*big.Int{NewInt(6), NewInt(-3), NewInt(0)}
	var e = make([]*bn256.G2, 0, 3)
	for i, val := range xi {
		e = append(e, new(bn256.G2).ScalarMult(val, c[i]))
	}

	var res = e[0]
	for i := 1; i < 3; i++ {
		res = new(bn256.G2).Add(e[i], res)
	}

	return new(bn256.G2).ScalarMult(res, inverse4)
}

func TestProving(_ *testing.T) {
	params := Setup(l1, l2, r2, o2, 3)
	proof := MakeProof(params, big1, big2, big2, big2, h)

	if err := VerifyProof(params, proof); err != nil {
		panic(err)
	}
}
