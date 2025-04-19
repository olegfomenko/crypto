// Package parcs
// Copyright 2025 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parcs

import (
	"crypto/rand"
	"math/big"
)

var one = big.NewInt(1)
var two = big.NewInt(2)

func f(x, c, n *big.Int) *big.Int {
	res := new(big.Int).Mul(x, x)
	res.Add(res, c)
	res.Mod(res, n)
	return res
}

func Pollard(n *big.Int) *big.Int {
	if n.Bit(0) == 0 {
		return two
	}

	for {
		// Constant for F
		c, err := rand.Int(rand.Reader, n)
		if err != nil {
			panic(err)
		}

		// Same rand value but will be used in different calculations
		x, err := rand.Int(rand.Reader, n)
		if err != nil {
			panic(err)
		}

		y := new(big.Int).Set(x)

		d := big.NewInt(1)

		for d.Cmp(one) == 0 {
			x = f(x, c, n)          // +1 step
			y = f(f(y, c, n), c, n) // +2 step

			diff := new(big.Int).Sub(x, y)
			diff.Abs(diff)

			d = new(big.Int).GCD(nil, nil, n, diff)

			//fmt.Println(n, x, y, d)
		}

		if d.Cmp(n) != 0 {
			return d
		}
	}

	panic("Never happens")
}
