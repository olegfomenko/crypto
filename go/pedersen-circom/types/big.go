// Package types
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package types

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

type BigInt struct {
	*big.Int
}

func Wrap(v *big.Int) *BigInt {
	return &BigInt{v}
}

func MustFromString(s string) *BigInt {
	return &BigInt{MustBigInt(s)}
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" || string(p) == "nil" {
		return nil
	}

	b.Int = MustBigInt(string(p))
	return nil
}

func (b *BigInt) Add(x *BigInt, y *BigInt) *BigInt {
	b.Int = new(big.Int).Mod(new(big.Int).Add(x.Int, y.Int), babyjub.Order)
	return b
}

func (b *BigInt) Mul(x *BigInt, y *BigInt) *BigInt {
	b.Int = new(big.Int).Mod(new(big.Int).Mul(x.Int, y.Int), babyjub.Order)
	return b
}

func (b *BigInt) Minus(val *BigInt) *BigInt {
	b.Int = new(big.Int).Mod(new(big.Int).Mul(val.Int, big.NewInt(-1)), babyjub.Order)
	return b
}

func MustBigInt(s string) *big.Int {
	res, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("failed to parse str to big.Int, str=%s", s))
	}
	return res
}
