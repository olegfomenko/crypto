// Package fft
// Copyright 2024 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package fft

import "math/big"

func FFT(p []*big.Int, domain []*big.Int, m *big.Int) []*big.Int {
	if len(p) == 1 {
		return p
	}

	l := FFT(evens(p), evens(domain), m)
	r := FFT(odds(p), evens(domain), m)

	res := make([]*big.Int, len(p))
	for i := 0; i < len(res)/2; i++ {
		rshift := new(big.Int).Mod(new(big.Int).Mul(domain[i], r[i]), m)
		res[i] = new(big.Int).Mod(new(big.Int).Add(l[i], rshift), m)
		res[i+len(res)/2] = new(big.Int).Mod(new(big.Int).Sub(l[i], rshift), m)
	}
	return res
}

func FFTInverse(p []*big.Int, domain []*big.Int, m *big.Int) []*big.Int {
	vals := FFT(p, domain, m)
	res := make([]*big.Int, len(p))
	ninv := new(big.Int).ModInverse(big.NewInt(int64(len(p))), m)

	res[0] = new(big.Int).Mod(new(big.Int).Mul(vals[0], ninv), m)

	for i := 1; i < len(res); i++ {
		res[i] = new(big.Int).Mod(new(big.Int).Mul(vals[len(vals)-i], ninv), m)
	}

	return res
}

func odds(p []*big.Int) []*big.Int {
	res := make([]*big.Int, 0, len(p)/2)
	for i := 1; i < len(p); i += 2 {
		res = append(res, new(big.Int).Set(p[i]))
	}
	return res
}

func evens(p []*big.Int) []*big.Int {
	res := make([]*big.Int, 0, len(p)/2)
	for i := 0; i < len(p); i += 2 {
		res = append(res, new(big.Int).Set(p[i]))
	}
	return res
}
