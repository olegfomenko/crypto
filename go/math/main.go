// Package math
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package math

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

func GCD(a, b *big.Int) *big.Int {
	if b.Cmp(big.NewInt(0)) == 0 {
		return a
	}

	return GCD(b, new(big.Int).Mod(a, b))
}

func Mu(n *big.Int) int {
	if n.Cmp(big.NewInt(1)) == 0 {
		return 1
	}

	cnt := 0

	for i := big.NewInt(2); new(big.Int).Mul(i, i).Cmp(n) <= 0; i = i.Add(i, big.NewInt(1)) {
		if new(big.Int).Mod(n, i).Cmp(big.NewInt(0)) == 0 {
			ai := 0
			for new(big.Int).Mod(n, i).Cmp(big.NewInt(0)) == 0 {
				n.Div(n, i)
				ai++
			}

			if ai > 1 {
				return 0
			}

			cnt++
		}
	}

	if n.Cmp(big.NewInt(1)) > 0 {
		cnt++
	}

	if cnt%2 == 0 {
		return 1
	}

	return -1
}

func Phi(n *big.Int) *big.Int {
	result := new(big.Int).Set(n)

	for i := big.NewInt(2); new(big.Int).Mul(i, i).Cmp(n) <= 0; i = i.Add(i, big.NewInt(1)) {
		if new(big.Int).Mod(n, i).Cmp(big.NewInt(0)) == 0 {
			for new(big.Int).Mod(n, i).Cmp(big.NewInt(0)) == 0 {
				n.Div(n, i)
			}

			result = new(big.Int).Sub(result, new(big.Int).Div(result, i))
		}
	}

	if n.Cmp(big.NewInt(1)) > 0 {
		result = new(big.Int).Sub(result, new(big.Int).Div(result, n))
	}

	return result
}

// FindSquareRoot uses Cipolla algorithm to solve x^2 = a (mod P)
// More information: https://en.wikipedia.org/wiki/Cipolla%27s_algorithm
func FindSquareRoot(n *big.Int, p *big.Int) (*big.Int, error) {
	n = new(big.Int).Set(n)
	p = new(big.Int).Set(p)

	for {
		a, err := rand.Int(rand.Reader, p)
		if err != nil {
			return nil, err
		}

		l, err := Legendre(sub(mul(a, a, p), n, p), p)
		if err != nil {
			return nil, err
		}

		if l.Cmp(big.NewInt(1)) == 0 {
			continue
		}

		i := new(big.Int).Sub(new(big.Int).Mul(a, a), n)

		f := &F2{
			x:      a,
			y:      big.NewInt(1),
			i:      i,
			modulo: p,
		}

		exp := new(big.Int).Div(new(big.Int).Add(p, big.NewInt(1)), big.NewInt(2))
		f = powF2(f, exp)
		return f.x, nil
	}
}

// TestPrime implements the Solovay–Strassen test for prime numbers
// More info: https://en.wikipedia.org/wiki/Solovay%E2%80%93Strassen_primality_test
func TestPrime(n *big.Int) (bool, error) {
	const k = 20
	n.Abs(n)

	if n.Cmp(big.NewInt(2)) == 0 {
		return true, nil
	}

	if new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 || n.Cmp(big.NewInt(1)) == 0 {
		return false, nil
	}

	for i := 0; i < k; i++ {
		a, err := rand.Int(rand.Reader, n)
		if err != nil {
			return false, err
		}

		if a.Cmp(big.NewInt(2)) <= 0 {
			continue
		}

		if GCD(a, n).Cmp(big.NewInt(1)) != 0 {
			return false, nil
		}

		j, err := Jacobi(a, n)
		if err != nil {
			return false, err
		}

		// a ^ ((n-1)/2)
		t := new(big.Int).Exp(
			a,
			new(big.Int).Div(new(big.Int).Sub(n, big.NewInt(1)), big.NewInt(2)),
			n,
		)

		j = new(big.Int).Mod(j.Add(j, n), n)

		if j.Cmp(t) != 0 {
			fmt.Println("failed", i, j, t, n, a)
			return false, nil
		}
	}

	return true, nil
}

func Legendre(a, p *big.Int) (*big.Int, error) {
	a = new(big.Int).Set(a)
	p = new(big.Int).Set(p)

	if p.Cmp(big.NewInt(2)) == 0 {
		return nil, errors.New("p should be prime > 2")
	}

	ok, err := TestPrime(p)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("invalid p: should be prime")
	}

	return Jacobi(a, p)
}

// Jacobi method calculates Jacobi symbol. About: https://en.wikipedia.org/wiki/Jacobi_symbol
// Definition:
// Jacobi(a, p) = Legendre(a,p1)*Legendre(a,p2)...*Legendre(a,pn)
func Jacobi(a *big.Int, p *big.Int) (r *big.Int, err error) {
	a = new(big.Int).Set(a)
	p = new(big.Int).Set(p)

	if new(big.Int).Mod(p, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("invalid p: should be 2k+1, k > 0")
	}

	if GCD(a, p).Cmp(big.NewInt(1)) != 0 {
		return big.NewInt(0), nil
	}

	r = big.NewInt(1)

	// Leverages on multiplicativity
	if a.Cmp(big.NewInt(0)) < 0 {
		a.Neg(a)
		// Check Jacobi(-1, p) definition
		if new(big.Int).Mod(p, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
			r.Neg(r)
		}
	}

	for {

		t := 0
		for {
			if new(big.Int).Mod(a, big.NewInt(2)).Cmp(big.NewInt(0)) != 0 {
				break
			}

			t++
			a.Div(a, big.NewInt(2))
		}

		// Check Jacobi(2, p) definition
		if t%2 == 1 {
			if new(big.Int).Mod(p, big.NewInt(8)).Cmp(big.NewInt(3)) == 0 ||
				new(big.Int).Mod(p, big.NewInt(8)).Cmp(big.NewInt(5)) == 0 {
				r.Neg(r)
			}
		}

		// law of quadratic reciprocity
		if new(big.Int).Mod(a, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 &&
			new(big.Int).Mod(p, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
			r.Neg(r)
		}

		c := a
		a = new(big.Int).Mod(p, c)
		p = c

		if a.Cmp(big.NewInt(0)) == 0 {
			return
		}
	}
}
