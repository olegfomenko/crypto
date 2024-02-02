// Package bp
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bp

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

type RangeProofSetup struct {
	*InnerArgumentSetup
	G *bn256.G1
	H *bn256.G1
}

type InnerArgumentSetup struct {
	n int
	g []*bn256.G1
	h []*bn256.G1
	u *bn256.G1
}

func newRangeProofSetup(n int) *RangeProofSetup {
	return &RangeProofSetup{
		InnerArgumentSetup: newInnerArgumentSetup(n),
		G:                  points(1)[0],
		H:                  points(1)[0],
	}
}

func newInnerArgumentSetup(n int) *InnerArgumentSetup {
	g := points(n)
	h := points(n)
	u := points(1)

	return &InnerArgumentSetup{
		n: n,
		g: g,
		h: h,
		u: u[0],
	}
}

//
func TestRangeProofWithInnerProduct(t *testing.T) {
	// a bit size of commitment value
	const n = 4
	// Public values (generators, etc.)
	setup := newRangeProofSetup(n)

	// commitment randomness
	prv := values(1)[0]

	// commitment value
	v := big.NewInt(11)

	// 1011 - bit representation of dec 11
	al := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(1)}
	// ar = al - 1^n
	ar := vectorSub(al, ones(n))

	// commitment to `v` with randomness `prv`
	V := com(setup.G, setup.H, v, prv)

	// commitment to al and ar randomness
	alpha := values(1)[0]

	// commitment to al and ar
	A := new(bn256.G1).Add(vecCom(setup.g, setup.h, al, ar), new(bn256.G1).ScalarMult(setup.H, alpha))

	// blinding vectors
	sl := values(n)
	sr := values(n)

	// commitment to sl and sr randomness
	ro := values(1)[0]
	S := new(bn256.G1).Add(vecCom(setup.g, setup.h, sl, sr), new(bn256.G1).ScalarMult(setup.H, ro))

	// Send A, S, V to verifier

	// Verifier generates challenges and sends it to prover
	y := values(1)[0]
	z := values(1)[0]

	yn := ntharr(y, n)
	z2 := mul(z, z)
	z3 := mul(z2, z)

	onen := ones(n)
	twon := ntharr(big.NewInt(2), n)

	t1 := add(
		vectorMul(hadamardMul(yn, sr), vectorSub(al, vectorMulOnScalar(onen, z))),
		vectorMul(sl, vectorAdd(vectorMulOnScalar(twon, z2), hadamardMul(yn, vectorAdd(ar, vectorMulOnScalar(onen, z))))),
	)

	t2 := vectorMul(hadamardMul(yn, sr), sl)

	// commitments to t1 and t2 randomness
	tau1 := values(1)[0]
	tau2 := values(1)[0]

	T1 := com(setup.G, setup.H, t1, tau1)
	T2 := com(setup.G, setup.H, t2, tau2)

	// Send T1, T2 to verifier

	// Verifier generates challenges and sends it to prover
	x := values(1)[0]

	x2 := mul(x, x)

	l := vectorAdd(vectorSub(al, vectorMulOnScalar(onen, z)), vectorMulOnScalar(sl, x))
	r := vectorAdd(
		hadamardMul(yn,
			vectorAdd(ar,
				vectorAdd(vectorMulOnScalar(onen, z), vectorMulOnScalar(sr, x)),
			),
		),
		vectorMulOnScalar(twon, z2),
	)

	tx := vectorMul(l, r)
	taux := add(mul(tau2, x2), add(mul(tau1, x), mul(z2, prv)))
	nu := add(alpha, mul(ro, x))

	// Prover sends tx, taux, nu to verifier

	// Verifier calculates:

	// 1. h1 := h^(y^-1)

	yinvn := invntharr(y, n) // [1, y^-1, y^-2, ... , y^-n+1]

	h1 := make([]*bn256.G1, n)
	for i := range h1 {
		h1[i] = new(bn256.G1).ScalarMult(setup.h[i], yinvn[i])
	}

	// 2. check that tx = t(x) = t0 + t1*x +t2*x^2

	deltayz := sub(mul(sub(z, z2), vectorMul(onen, yn)), mul(z3, vectorMul(onen, twon)))

	c1 := com(setup.G, setup.H, tx, taux)

	c2 := new(bn256.G1).ScalarMult(V, z2)
	c2.Add(c2, new(bn256.G1).ScalarMult(setup.G, deltayz))
	c2.Add(c2, new(bn256.G1).ScalarMult(T1, x))
	c2.Add(c2, new(bn256.G1).ScalarMult(T2, x2))

	if !bytes.Equal(c1.Marshal(), c2.Marshal()) {
		panic("Failed: tx ?= t0 + t1*x +t2*x^2")
	}

	// 3. compute commitment to l, r

	P := new(bn256.G1).Add(A, new(bn256.G1).ScalarMult(S, x))
	P.Add(P, vectorPointScalarMul(setup.g, vectorMulOnScalar(onen, sub(big.NewInt(0), z))))
	P.Add(P, vectorPointScalarMul(h1, vectorAdd(vectorMulOnScalar(yn, z), vectorMulOnScalar(twon, z2))))

	// P = h^nu * g^l * g^r
	// For inner product use: P* h^-nu * u^t

	P.Add(P, new(bn256.G1).ScalarMult(setup.H, sub(big.NewInt(0), nu)))
	P.Add(P, new(bn256.G1).ScalarMult(setup.u, tx))

	// 4. Run inner product proof on (n, g, h1, P* h^-nu * u^t, l, r)
	protocol2(n, setup.g, h1, setup.u, P, l, r)
}

func TestRangeProof(t *testing.T) {
	// a bit size of commitment value
	const n = 4
	// Public values (generators, etc.)
	setup := newRangeProofSetup(n)

	// commitment randomness
	prv := values(1)[0]

	// commitment value
	v := big.NewInt(11)

	// 1011 - bit representation of dec 11
	al := []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(1)}
	// ar = al - 1^n
	ar := vectorSub(al, ones(n))

	// commitment to `v` with randomness `prv`
	V := com(setup.G, setup.H, v, prv)

	// commitment to al and ar randomness
	alpha := values(1)[0]

	// commitment to al and ar
	A := new(bn256.G1).Add(vecCom(setup.g, setup.h, al, ar), new(bn256.G1).ScalarMult(setup.H, alpha))

	// blinding vectors
	sl := values(n)
	sr := values(n)

	// commitment to sl and sr randomness
	ro := values(1)[0]
	S := new(bn256.G1).Add(vecCom(setup.g, setup.h, sl, sr), new(bn256.G1).ScalarMult(setup.H, ro))

	// Send A, S, V tp verifier

	// Verifier generates challenges and sends it to prover
	y := values(1)[0]
	z := values(1)[0]

	yn := ntharr(y, n)
	z2 := mul(z, z)
	z3 := mul(z2, z)

	onen := ones(n)
	twon := ntharr(big.NewInt(2), n)

	t1 := add(
		vectorMul(hadamardMul(yn, sr), vectorSub(al, vectorMulOnScalar(onen, z))),
		vectorMul(sl, vectorAdd(vectorMulOnScalar(twon, z2), hadamardMul(yn, vectorAdd(ar, vectorMulOnScalar(onen, z))))),
	)

	t2 := vectorMul(hadamardMul(yn, sr), sl)

	// commitments to t1 and t2 randomness
	tau1 := values(1)[0]
	tau2 := values(1)[0]

	T1 := com(setup.G, setup.H, t1, tau1)
	T2 := com(setup.G, setup.H, t2, tau2)

	// Send T1, T2 to verifier

	// Verifier generates challenges and sends it to prover
	x := values(1)[0]

	x2 := mul(x, x)

	l := vectorAdd(vectorSub(al, vectorMulOnScalar(onen, z)), vectorMulOnScalar(sl, x))
	r := vectorAdd(
		hadamardMul(yn,
			vectorAdd(ar,
				vectorAdd(vectorMulOnScalar(onen, z), vectorMulOnScalar(sr, x)),
			),
		),
		vectorMulOnScalar(twon, z2),
	)

	tx := vectorMul(l, r)
	taux := add(mul(tau2, x2), add(mul(tau1, x), mul(z2, prv)))
	nu := add(alpha, mul(ro, x))

	// Prover sends l, r, tx, taux, nu to verifier

	// Verifier calculates:

	// 1. h1 := h^(y^-1)

	yinvn := invntharr(y, n) // [1, y^-1, y^-2, ... , y^-n+1]

	h1 := make([]*bn256.G1, n)
	for i := range h1 {
		h1[i] = new(bn256.G1).ScalarMult(setup.h[i], yinvn[i])
	}

	// 2. check that tx = t(x) = t0 + t1*x +t2*x^2

	deltayz := sub(mul(sub(z, z2), vectorMul(onen, yn)), mul(z3, vectorMul(onen, twon)))

	c1 := com(setup.G, setup.H, tx, taux)

	c2 := new(bn256.G1).ScalarMult(V, z2)
	c2.Add(c2, new(bn256.G1).ScalarMult(setup.G, deltayz))
	c2.Add(c2, new(bn256.G1).ScalarMult(T1, x))
	c2.Add(c2, new(bn256.G1).ScalarMult(T2, x2))

	if !bytes.Equal(c1.Marshal(), c2.Marshal()) {
		panic("Failed: tx ?= t0 + t1*x +t2*x^2")
	}

	// 3. compute commitment to l, r

	P := new(bn256.G1).Add(A, new(bn256.G1).ScalarMult(S, x))
	P.Add(P, vectorPointScalarMul(setup.g, vectorMulOnScalar(onen, sub(big.NewInt(0), z))))
	P.Add(P, vectorPointScalarMul(h1, vectorAdd(vectorMulOnScalar(yn, z), vectorMulOnScalar(twon, z2))))

	// 4. Check that l,r are valid

	P1 := vecCom(setup.g, h1, l, r)
	P1.Add(P1, new(bn256.G1).ScalarMult(setup.H, nu))

	if !bytes.Equal(P.Marshal(), P1.Marshal()) {
		panic("Failed: l, r are invalid")
	}

	// 5. Check t =<l,r>
	if tx.Cmp(vectorMul(l, r)) != 0 {
		panic("Failed: t ?= <l, r>")
	}
}

func TestInnerArgumentProofRecursive(t *testing.T) {
	// public values
	const n = 4
	setup := newInnerArgumentSetup(n)

	// private values
	a := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(10), big.NewInt(1)}
	b := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(10)}

	// public commitment P: p = g^a+h^b+u^<a, b>
	p := productCom(setup.g, setup.h, setup.u, a, b)

	// Run proof system
	protocol2(setup.n, setup.g, setup.h, setup.u, p, a, b)
}

func protocol2(n int, g, h []*bn256.G1, u *bn256.G1, p *bn256.G1, a, b []*big.Int) {
	fmt.Println("Running protocol2() in n =", n)

	if n == 1 {
		// send a, b to verifier

		// then verifier computes
		p1 := new(bn256.G1).Add(com(g[0], h[0], a[0], b[0]), new(bn256.G1).ScalarMult(u, mul(a[0], b[0])))

		// Verifier perform check...
		if !bytes.Equal(p1.Marshal(), p.Marshal()) {
			panic("Proof is invalid")
		}

		// Successful
		fmt.Println("Proof is valid")
		return
	}

	if n%2 != 0 {
		panic("invalid n")
	}

	n1 := n / 2

	// Same as we are using H in non-recursive version
	cl := vectorMul(a[:n1], b[n1:])
	cr := vectorMul(a[n1:], b[:n1])

	L := new(bn256.G1).Add(vecCom(g[n1:], h[:n1], a[:n1], b[n1:]), new(bn256.G1).ScalarMult(u, cl))
	R := new(bn256.G1).Add(vecCom(g[:n1], h[n1:], a[n1:], b[:n1]), new(bn256.G1).ScalarMult(u, cr))

	// Send L, R to verifier

	// Verifier generates challenge and sends to prover
	x := values(1)[0]

	xinv := new(big.Int).ModInverse(x, bn256.Order)
	x2 := mul(x, x)
	x2inv := new(big.Int).ModInverse(x2, bn256.Order)

	// Both verifier and prover computes new public values

	g1 := hadamardPointMul(vectorPointMulOnScalar(g[:n1], xinv), vectorPointMulOnScalar(g[n1:], x))
	h1 := hadamardPointMul(vectorPointMulOnScalar(h[:n1], x), vectorPointMulOnScalar(h[n1:], xinv))

	p1 := new(bn256.G1).Add(new(bn256.G1).ScalarMult(L, x2), p)
	p1.Add(p1, new(bn256.G1).ScalarMult(R, x2inv)) // p1 := L^x2*P*R^x-2

	// Prover computes new a,b
	a1 := vectorAdd(vectorMulOnScalar(a[:n1], x), vectorMulOnScalar(a[n1:], xinv))
	b1 := vectorAdd(vectorMulOnScalar(b[n1:], x), vectorMulOnScalar(b[:n1], xinv))

	protocol2(n1, g1, h1, u, p1, a1, b1)
}

func TestInnerArgumentProof(t *testing.T) {
	// public values
	const n = 4
	const n1 = n / 2
	g := points(n)
	h := points(n)
	u := points(1)

	// private values
	a := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(10), big.NewInt(1)}
	b := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(10)}
	c := vectorMul(a, b) // a * b = 4 + 5 + 20 + 10 = 39

	L := InnerProductH(n, g, h, u[0], zeros(n1), a[:n1], b[n1:], zeros(n1), vectorMul(a[:n1], b[n1:]))
	R := InnerProductH(n, g, h, u[0], a[n1:], zeros(n1), zeros(n1), b[:n1], vectorMul(a[n1:], b[:n1]))

	P := InnerProductH(n, g, h, u[0], a[:n1], a[n1:], b[:n1], b[n1:], c) // main commitment: p = g^a+h^b+u^<a, b>

	// Send L, R, P to verifier

	// Verifier generates challenge and sends to prover
	x, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}

	// Prover calculates...

	xinv := new(big.Int).ModInverse(x, bn256.Order)

	a1 := vectorAdd(vectorMulOnScalar(a[:n1], x), vectorMulOnScalar(a[n1:], xinv))
	b1 := vectorAdd(vectorMulOnScalar(b[n1:], x), vectorMulOnScalar(b[:n1], xinv))

	// Send a1, b1 to verifier

	x2 := mul(x, x)

	lx2 := new(bn256.G1).ScalarMult(L, x2)
	rx2 := new(bn256.G1).ScalarMult(R, new(big.Int).ModInverse(x2, bn256.Order))

	p1 := new(bn256.G1).Add(lx2, P)
	p1.Add(p1, rx2)

	pcheck := InnerProductH(n, g, h, u[0], vectorMulOnScalar(a1, xinv), vectorMulOnScalar(a1, x), vectorMulOnScalar(b1, x), vectorMulOnScalar(b1, xinv), vectorMul(a1, b1))

	// Verifier perform check...
	if !bytes.Equal(p1.Marshal(), pcheck.Marshal()) {
		panic("check fails")
	}
}
