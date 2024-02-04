// Package bp
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bp

import (
	"crypto/rand"
	"math/big"

	"github.com/cloudflare/bn256"
	"github.com/iden3/go-iden3-crypto/keccak256"
)

func InnerProductH(n int, g []*bn256.G1, h []*bn256.G1, u *bn256.G1, a, a1, b, b1 []*big.Int, c *big.Int) *bn256.G1 {
	if n%2 != 0 {
		panic("invalid n")
	}

	n1 := n / 2
	res := vectorPointScalarMul(g[:n1], a)
	res.Add(res, vectorPointScalarMul(g[n1:], a1))
	res.Add(res, vectorPointScalarMul(h[:n1], b))
	res.Add(res, vectorPointScalarMul(h[n1:], b1))
	res.Add(res, new(bn256.G1).ScalarMult(u, c))
	return res
}

func invntharr(val *big.Int, n int) []*big.Int {
	res := make([]*big.Int, n)

	res[0] = big.NewInt(1)

	inv := new(big.Int).ModInverse(val, bn256.Order)

	for i := 1; i < n; i++ {
		res[i] = mul(res[i-1], inv)
	}
	return res
}

func ntharr(val *big.Int, n int) []*big.Int {
	res := make([]*big.Int, n)
	res[0] = big.NewInt(1)
	for i := 1; i < n; i++ {
		res[i] = mul(res[i-1], val)
	}
	return res
}

func values(n int) []*big.Int {
	res := make([]*big.Int, n)
	var err error

	for i := range res {
		res[i], err = rand.Int(rand.Reader, bn256.Order)
		if err != nil {
			panic(err)
		}
	}

	return res
}

func ones(n int) []*big.Int {
	res := make([]*big.Int, n)
	for i := range res {
		res[i] = big.NewInt(1)
	}
	return res
}

func zeros(n int) []*big.Int {
	res := make([]*big.Int, n)
	for i := range res {
		res[i] = big.NewInt(0)
	}
	return res
}

func points(n int) []*bn256.G1 {
	res := make([]*bn256.G1, n)
	var err error
	for i := range res {
		if _, res[i], err = bn256.RandomG1(rand.Reader); err != nil {
			panic(err)
		}
	}
	return res
}

func com(g, h *bn256.G1, a, b *big.Int) *bn256.G1 {
	return new(bn256.G1).Add(new(bn256.G1).ScalarMult(g, a), new(bn256.G1).ScalarMult(h, b))
}

func vecCom(g, h []*bn256.G1, a, b []*big.Int) *bn256.G1 {
	return new(bn256.G1).Add(vectorPointScalarMul(g, a), vectorPointScalarMul(h, b))
}

func productCom(g []*bn256.G1, h []*bn256.G1, u *bn256.G1, a, b []*big.Int) *bn256.G1 {
	p := vecCom(g, h, a, b)
	p.Add(p, new(bn256.G1).ScalarMult(u, vectorMul(a, b)))
	return p
}

func vectorAdd(a []*big.Int, b []*big.Int) []*big.Int {
	if len(b) != len(a) {
		panic("invalid length")
	}

	res := make([]*big.Int, len(a))
	for i := 0; i < len(res); i++ {
		res[i] = add(a[i], b[i])
	}

	return res
}

func vectorSub(a []*big.Int, b []*big.Int) []*big.Int {
	if len(b) != len(a) {
		panic("invalid length")
	}

	res := make([]*big.Int, len(a))
	for i := 0; i < len(res); i++ {
		res[i] = sub(a[i], b[i])
	}

	return res
}

func vectorMulOnScalar(a []*big.Int, c *big.Int) []*big.Int {
	res := make([]*big.Int, len(a))
	for i := range res {
		res[i] = mul(a[i], c)
	}
	return res
}

func vectorMul(a []*big.Int, b []*big.Int) *big.Int {
	if len(b) != len(a) {
		panic("invalid length")
	}

	res := big.NewInt(0)
	for i := 0; i < len(a); i++ {
		res = add(res, mul(a[i], b[i]))
	}
	return res
}

func hadamardMul(a, b []*big.Int) []*big.Int {
	if len(b) != len(a) {
		panic("invalid length")
	}

	res := make([]*big.Int, len(a))
	for i := range res {
		res[i] = mul(a[i], b[i])
	}

	return res
}

func vectorPointScalarMul(g []*bn256.G1, a []*big.Int) *bn256.G1 {
	if len(g) != len(a) {
		panic("invalid length for scalar mul")
	}

	res := new(bn256.G1).ScalarMult(g[0], a[0])
	for i := 1; i < len(g); i++ {
		res.Add(res, new(bn256.G1).ScalarMult(g[i], a[i]))
	}
	return res
}

func vectorPointMulOnScalar(g []*bn256.G1, a *big.Int) []*bn256.G1 {
	res := make([]*bn256.G1, len(g))
	for i := range res {
		res[i] = new(bn256.G1).ScalarMult(g[i], a)
	}
	return res
}

func hadamardPointMul(a, b []*bn256.G1) []*bn256.G1 {
	if len(a) != len(b) {
		panic("invalid length for scalar mul")
	}

	res := make([]*bn256.G1, len(a))
	for i := range res {
		res[i] = new(bn256.G1).Add(a[i], b[i])
	}
	return res
}

func add(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(x, y), bn256.Order)
}

func sub(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Sub(x, y), bn256.Order)
}

func mul(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, y), bn256.Order)
}

func scalarTo32Byte(s *big.Int) []byte {
	return to32Bytes(s.Bytes())
}

func to32Bytes(arr []byte) []byte {
	if len(arr) >= 32 {
		return arr[:32]
	}

	res := make([]byte, 32-len(arr))
	return append(res, arr...)
}

func hash(scalars []*big.Int, points []*bn256.G1) *big.Int {
	var data []byte
	for _, s := range scalars {
		data = append(data, scalarTo32Byte(s)...)
	}

	for _, p := range points {
		data = append(data, p.Marshal()...)
	}

	return new(big.Int).Mod(new(big.Int).SetBytes(keccak256.Hash(data)), bn256.Order)
}

func toBits(v *big.Int, n int) []*big.Int {
	res := make([]*big.Int, 0, n)
	str := v.Text(2)

	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '1' {
			res = append(res, big.NewInt(1))
			continue
		}
		res = append(res, big.NewInt(0))
	}

	for len(res) < n {
		res = append(res, big.NewInt(0))
	}
	return res[:n]
}
