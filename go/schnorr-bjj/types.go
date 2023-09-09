// Package schnorr_bjj
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package schnorr_bjj

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

var G = babyjub.B8

type SchnorrSignature struct {
	R *babyjub.Point
	S *big.Int
}
