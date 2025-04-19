// Package parcs
// Copyright 2025 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parcs

import (
	"fmt"
	"math/big"
	"testing"
)

func TestParcsPollard_100(t *testing.T) {
	p := NewParcsPollard(10)
	res := p.Run(big.NewInt(100))
	fmt.Println(res)
}

func TestParcsPollard_8051(t *testing.T) {
	p := NewParcsPollard(10)
	res := p.Run(big.NewInt(8051))
	fmt.Println(res)
}

// N=11×17×23×37×53×101×113×127×149×191=347912642190594349.
func TestParcsPollard_347912642190594349(t *testing.T) {
	p := NewParcsPollard(5)
	res := p.Run(big.NewInt(347912642190594349))
	fmt.Println(res)
}
