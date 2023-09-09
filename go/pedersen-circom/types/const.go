// Package types
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package types

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	xH, _ = new(big.Int).SetString("15334330715717027115948243110556436026028216985345384579806128223314358448928", 10)
	yH, _ = new(big.Int).SetString("14640338696677432581567520324796424956409796398271990973432884194068091890885", 10)

	G = babyjub.B8
	H = &babyjub.Point{X: xH, Y: yH}
)
