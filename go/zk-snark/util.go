// Package zk_snark
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package zk_snark

import (
	"crypto/rand"
	"math/big"
)

func GetRandomWithMax(max *big.Int) *big.Int {
	val, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	return val
}
