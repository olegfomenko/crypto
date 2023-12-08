// Package math
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package math

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestPhi(t *testing.T) {
	fmt.Println(Phi(big.NewInt(31))) // 30
	fmt.Println(Phi(big.NewInt(10))) // 4
	fmt.Println(Phi(big.NewInt(7)))  // 7
}

func TestFindSquareRoot(t *testing.T) {
	fmt.Println(FindSquareRoot(big.NewInt(10), big.NewInt(13)))    // 6 or 7
	fmt.Println(FindSquareRoot(big.NewInt(362), big.NewInt(7919))) // 7828 or 91
}

func TestJacobi(t *testing.T) {
	tests := []struct {
		a   *big.Int
		p   *big.Int
		res *big.Int
	}{
		{
			a:   big.NewInt(7),
			p:   big.NewInt(35),
			res: big.NewInt(0),
		},
		{
			a:   big.NewInt(2),
			p:   big.NewInt(41),
			res: big.NewInt(1),
		},
		{
			a:   big.NewInt(21),
			p:   big.NewInt(9),
			res: big.NewInt(0),
		},
		{
			a:   big.NewInt(8),
			p:   big.NewInt(13),
			res: big.NewInt(-1),
		},
		{
			a:   big.NewInt(4),
			p:   big.NewInt(55),
			res: big.NewInt(1),
		},
		{
			a:   big.NewInt(9),
			p:   big.NewInt(37),
			res: big.NewInt(1),
		},
	}

	for i, t := range tests {
		fmt.Printf("Running %d\n", i)
		res, err := Jacobi(t.a, t.p)
		if err != nil {
			panic(err)
		}

		if res.Cmp(t.res) != 0 {
			panic(fmt.Sprintf("test cast %d failed", i))
		}
		fmt.Printf("Finished %d\n", i)
	}
}

func TestTestPrimes(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Printf("Running %d prime check\n", i)
		prime, err := rand.Prime(rand.Reader, 256)
		if err != nil {
			panic(err)
		}

		ok, err := TestPrime(prime)
		if err != nil {
			panic(err)
		}

		if !ok {
			fmt.Println(prime.String())
			panic("verification failed on prime number")
		}
	}
}

func TestLegendre(t *testing.T) {
	fmt.Println(Legendre(big.NewInt(219), big.NewInt(383))) // Should be 1
}
