// Package math
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package math

import "math/big"

// F2 represents field extension for F(sqrt(i)): {x + y*sqrt(i) | x, y from F_modulo}
type F2 struct {
	x, y   *big.Int
	i      *big.Int
	modulo *big.Int
}

func addF2(val1, val2 *F2) *F2 {
	return &F2{
		add(val1.x, val2.x, val1.modulo),
		add(val1.y, val2.y, val1.modulo),
		val1.i,
		val1.modulo,
	}
}

func mulF2(val1, val2 *F2) *F2 {
	return &F2{
		x:      add(mul(val1.x, val2.x, val1.modulo), mul(mul(val1.y, val2.y, val1.modulo), val1.i, val1.modulo), val1.modulo),
		y:      add(mul(val1.x, val2.y, val1.modulo), mul(val1.y, val2.x, val1.modulo), val1.modulo),
		i:      val1.i,
		modulo: val1.modulo,
	}
}

func powF2(val *F2, exp *big.Int) *F2 {
	if exp.Cmp(big.NewInt(0)) == 0 {
		return &F2{
			big.NewInt(0),
			big.NewInt(0),
			val.i,
			val.modulo,
		}
	}

	if exp.Cmp(big.NewInt(1)) == 0 {
		return val
	}

	f := powF2(val, new(big.Int).Div(exp, big.NewInt(2)))
	f = mulF2(f, f)

	if new(big.Int).Mod(exp, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
		f = mulF2(f, val)
	}

	return f
}

func add(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(x, y), mod)
}

func sub(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Sub(x, y), mod)
}

func mul(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, y), mod)
}

func div(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, new(big.Int).ModInverse(y, mod)), mod)
}
