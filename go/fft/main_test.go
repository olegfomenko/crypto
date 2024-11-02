// Package fft
// Copyright 2024 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package fft

import (
	"fmt"
	"math/big"
	"testing"
)

func TestFFT(t *testing.T) {
	resFFT := FFT(
		[]*big.Int{big.NewInt(3), big.NewInt(1), big.NewInt(4), big.NewInt(1), big.NewInt(5), big.NewInt(9), big.NewInt(2), big.NewInt(6)},
		[]*big.Int{big.NewInt(1), big.NewInt(85), big.NewInt(148), big.NewInt(111), big.NewInt(336), big.NewInt(252), big.NewInt(189), big.NewInt(226)},
		big.NewInt(337),
	)

	fmt.Println(resFFT)

	resInv := FFTInverse(
		resFFT,
		[]*big.Int{big.NewInt(1), big.NewInt(85), big.NewInt(148), big.NewInt(111), big.NewInt(336), big.NewInt(252), big.NewInt(189), big.NewInt(226)},
		big.NewInt(337),
	)

	fmt.Println(resInv)
}
