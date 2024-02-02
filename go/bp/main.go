// Package bp
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package bp

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/cloudflare/bn256"
)

// BulletProofPublic represents public general information about range proof system.
// It can be used for all proofs.
type BulletProofPublic struct {
	*InnerArgumentPublic
	// N is a bit size of `v` value: 0 =< v <= 2^n - 1
	N int
	// Commitment base points
	G *bn256.G1
	H *bn256.G1
}

// NewBulletProofPublic generates new public data for given N.
func NewBulletProofPublic(n int) *BulletProofPublic {
	return &BulletProofPublic{
		InnerArgumentPublic: NewInnerArgumentPublic(n),
		N:                   n,
		G:                   points(1)[0],
		H:                   points(1)[0],
	}
}

// InnerArgumentPublic represents public general information for inner argument proof system.
// It can be used for all proofs.
type InnerArgumentPublic struct {
	// N is a size of input vectors
	N int
	// Commitment base points
	G []*bn256.G1
	H []*bn256.G1
	U *bn256.G1
}

// NewInnerArgumentPublic generates new public data for given N.
func NewInnerArgumentPublic(n int) *InnerArgumentPublic {
	return &InnerArgumentPublic{
		N: n,
		G: points(n),
		H: points(n),
		U: points(1)[0],
	}
}

// BulletProof represents range ZK proof and contains information about global parameters
// and all public data that is required to verify proof.
type BulletProof struct {
	*BulletProofPublic

	// Commitment to `v`: V = (g^v)*(h^prv)
	V *bn256.G1

	// Bulletproof values
	ACom  *bn256.G1
	SCom  *bn256.G1
	T1Com *bn256.G1
	T2Com *bn256.G1
	Tx    *big.Int
	TauX  *big.Int
	Nu    *big.Int

	// Inner product proof values
	L    []*bn256.G1
	R    []*bn256.G1
	A, B *big.Int
}

// Prove generates ZK range proof for given value `v` and randomness `prv` based on global parameters.
func (p *BulletProofPublic) Prove(v, prv *big.Int) (proof *BulletProof, err error) {
	onen := ones(p.N)

	al := toBits(v, p.N)
	ar := vectorSub(al, onen)

	alpha := values(1)[0]

	A := new(bn256.G1).Add(vecCom(p.InnerArgumentPublic.G, p.InnerArgumentPublic.H, al, ar), new(bn256.G1).ScalarMult(p.H, alpha))

	sl := values(p.N)
	sr := values(p.N)

	ro := values(1)[0]
	S := new(bn256.G1).Add(vecCom(p.InnerArgumentPublic.G, p.InnerArgumentPublic.H, sl, sr), new(bn256.G1).ScalarMult(p.H, ro))

	V := com(p.G, p.H, v, prv)

	proof = &BulletProof{
		BulletProofPublic: p,
		V:                 V,
		ACom:              A,
		SCom:              S,
	}

	// Using Fiat-Shamir
	y := hash([]*big.Int{big.NewInt(int64(p.N))}, []*bn256.G1{A, S, V})
	z := hash([]*big.Int{y}, []*bn256.G1{A, S})

	yn := ntharr(y, p.N)
	z2 := mul(z, z)

	twon := ntharr(big.NewInt(2), p.N)

	t1 := add(
		vectorMul(hadamardMul(yn, sr), vectorSub(al, vectorMulOnScalar(onen, z))),
		vectorMul(sl, vectorAdd(vectorMulOnScalar(twon, z2), hadamardMul(yn, vectorAdd(ar, vectorMulOnScalar(onen, z))))),
	)

	t2 := vectorMul(hadamardMul(yn, sr), sl)

	tau1 := values(1)[0]
	tau2 := values(1)[0]

	T1 := com(p.G, p.H, t1, tau1)
	T2 := com(p.G, p.H, t2, tau2)

	proof.T1Com = T1
	proof.T2Com = T2

	// Using Fiat-Shamir
	x := hash([]*big.Int{y, z}, []*bn256.G1{T1, T2})

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

	proof.Tx = tx
	proof.TauX = taux
	proof.Nu = nu

	yinvn := invntharr(y, p.N) // [1, y^-1, y^-2, ... , y^-n+1]

	h1 := make([]*bn256.G1, p.N)
	for i := range h1 {
		h1[i] = new(bn256.G1).ScalarMult(p.InnerArgumentPublic.H[i], yinvn[i])
	}

	public := &InnerArgumentPublic{
		N: p.N,
		G: p.InnerArgumentPublic.G,
		H: h1,
		U: p.InnerArgumentPublic.U,
	}

	innerProductProof, err := public.Prove(l, r)
	if err != nil {
		return nil, err
	}

	proof.A = innerProductProof.A
	proof.B = innerProductProof.B
	proof.L = innerProductProof.L
	proof.R = innerProductProof.R
	return proof, nil
}

// Verify verifies ZK range proof based on global parameters.
func (p *BulletProofPublic) Verify(proof *BulletProof) error {
	// Using Fiat-Shamir
	y := hash([]*big.Int{big.NewInt(int64(proof.N))}, []*bn256.G1{proof.ACom, proof.SCom, proof.V})
	z := hash([]*big.Int{y}, []*bn256.G1{proof.ACom, proof.SCom})

	yn := ntharr(y, p.N)
	z2 := mul(z, z)
	z3 := mul(z2, z)

	onen := ones(p.N)
	twon := ntharr(big.NewInt(2), p.N)

	// Using Fiat-Shamir
	x := hash([]*big.Int{y, z}, []*bn256.G1{proof.T1Com, proof.T2Com})

	x2 := mul(x, x)

	// Verifier calculates:

	// 1. h1 := h^(y^-1)

	yinvn := invntharr(y, p.N) // [1, y^-1, y^-2, ... , y^-n+1]

	h1 := make([]*bn256.G1, p.N)
	for i := range h1 {
		h1[i] = new(bn256.G1).ScalarMult(proof.InnerArgumentPublic.H[i], yinvn[i])
	}

	// 2. check that tx = t(x) = t0 + t1*x +t2*x^2

	deltayz := sub(mul(sub(z, z2), vectorMul(onen, yn)), mul(z3, vectorMul(onen, twon)))

	c1 := com(proof.G, proof.H, proof.Tx, proof.TauX)

	c2 := new(bn256.G1).ScalarMult(proof.V, z2)
	c2.Add(c2, new(bn256.G1).ScalarMult(proof.G, deltayz))
	c2.Add(c2, new(bn256.G1).ScalarMult(proof.T1Com, x))
	c2.Add(c2, new(bn256.G1).ScalarMult(proof.T2Com, x2))

	if !bytes.Equal(c1.Marshal(), c2.Marshal()) {
		return errors.New("failed: tx ?= t0 + t1*x +t2*x^2")
	}

	P := new(bn256.G1).Add(proof.ACom, new(bn256.G1).ScalarMult(proof.SCom, x))
	P.Add(P, vectorPointScalarMul(proof.InnerArgumentPublic.G, vectorMulOnScalar(onen, sub(big.NewInt(0), z))))
	P.Add(P, vectorPointScalarMul(h1, vectorAdd(vectorMulOnScalar(yn, z), vectorMulOnScalar(twon, z2))))

	// P = h^nu * g^l * g^r
	// For inner product use: P* h^-nu * u^t

	P.Add(P, new(bn256.G1).ScalarMult(proof.H, sub(big.NewInt(0), proof.Nu)))
	P.Add(P, new(bn256.G1).ScalarMult(proof.U, proof.Tx))

	public := &InnerArgumentPublic{
		N: p.N,
		G: p.InnerArgumentPublic.G,
		H: h1,
		U: p.InnerArgumentPublic.U,
	}

	return public.Verify(&InnerProductProof{
		InnerArgumentPublic: public,
		L:                   proof.L,
		R:                   proof.R,
		A:                   proof.A,
		B:                   proof.B,
		P:                   P,
	})
}

// InnerProductProof represents ZK proof of knowledge of committed product of two vectors (inner product argument proof).
// Contains information about global parameters and all public data that is required to verify proof.
type InnerProductProof struct {
	*InnerArgumentPublic
	L    []*bn256.G1
	R    []*bn256.G1
	A, B *big.Int
	// Vector commitment for inner product argument: V = (g^a)*(h^b)*(u^<a,b>)
	P *bn256.G1
}

// Prove generates ZK inner argument proof with logarithmic size for given vectors `a` and `b`.
func (p *InnerArgumentPublic) Prove(a []*big.Int, b []*big.Int) (*InnerProductProof, error) {
	return innerArgProof(p, a, b)
}

// Verify verifies ZK inner argument proof based on global parameters.
func (p *InnerArgumentPublic) Verify(proof *InnerProductProof) error {
	proof.InnerArgumentPublic = p
	return innerArgVerify(0, proof)
}

func innerArgVerify(deep int, proof *InnerProductProof) error {
	if proof.N == 1 {
		p := new(bn256.G1).Add(com(proof.G[0], proof.H[0], proof.A, proof.B), new(bn256.G1).ScalarMult(proof.U, mul(proof.A, proof.B)))

		// Verifier perform check...
		if !bytes.Equal(proof.P.Marshal(), p.Marshal()) {
			return errors.New("failed to verify inner argument proof")
		}

		return nil
	}

	if proof.N%2 != 0 {
		return errors.New("invalid n: should be 2^x")
	}

	n1 := proof.N / 2

	// Using Fiat-Shamir
	x := hash([]*big.Int{big.NewInt(int64(proof.N))}, []*bn256.G1{proof.P, proof.L[deep], proof.R[deep]})

	xinv := new(big.Int).ModInverse(x, bn256.Order)
	x2 := mul(x, x)
	x2inv := new(big.Int).ModInverse(x2, bn256.Order)

	// Both verifier and prover computes new public values

	g1 := hadamardPointMul(vectorPointMulOnScalar(proof.G[:n1], xinv), vectorPointMulOnScalar(proof.G[n1:], x))
	h1 := hadamardPointMul(vectorPointMulOnScalar(proof.H[:n1], x), vectorPointMulOnScalar(proof.H[n1:], xinv))

	p := new(bn256.G1).Add(new(bn256.G1).ScalarMult(proof.L[deep], x2), proof.P)
	p.Add(p, new(bn256.G1).ScalarMult(proof.R[deep], x2inv)) // p := L^x2*P*R^x-2

	return innerArgVerify(deep+1, &InnerProductProof{
		InnerArgumentPublic: &InnerArgumentPublic{
			N: n1,
			G: g1,
			H: h1,
			U: proof.U,
		},
		L: proof.L,
		R: proof.R,
		A: proof.A,
		B: proof.B,
		P: p,
	})
}

func innerArgProof(public *InnerArgumentPublic, a, b []*big.Int) (*InnerProductProof, error) {
	if public.N == 1 {
		// send a, b to verifier
		return &InnerProductProof{
			InnerArgumentPublic: public,
			A:                   a[0],
			B:                   b[0],
		}, nil
	}

	if public.N%2 != 0 {
		return nil, errors.New("invalid n: should be 2^x")
	}

	n1 := public.N / 2

	cl := vectorMul(a[:n1], b[n1:])
	cr := vectorMul(a[n1:], b[:n1])

	L := new(bn256.G1).Add(vecCom(public.G[n1:], public.H[:n1], a[:n1], b[n1:]), new(bn256.G1).ScalarMult(public.U, cl))
	R := new(bn256.G1).Add(vecCom(public.G[:n1], public.H[n1:], a[n1:], b[:n1]), new(bn256.G1).ScalarMult(public.U, cr))

	proof := &InnerProductProof{
		InnerArgumentPublic: public,
		L:                   []*bn256.G1{L},
		R:                   []*bn256.G1{R},
		P:                   productCom(public.G, public.H, public.U, a, b),
	}

	// Using Fiat-Shamir
	x := hash([]*big.Int{big.NewInt(int64(public.N))}, []*bn256.G1{proof.P, L, R})

	xinv := new(big.Int).ModInverse(x, bn256.Order)

	g1 := hadamardPointMul(vectorPointMulOnScalar(public.G[:n1], xinv), vectorPointMulOnScalar(public.G[n1:], x))
	h1 := hadamardPointMul(vectorPointMulOnScalar(public.H[:n1], x), vectorPointMulOnScalar(public.H[n1:], xinv))

	a1 := vectorAdd(vectorMulOnScalar(a[:n1], x), vectorMulOnScalar(a[n1:], xinv))
	b1 := vectorAdd(vectorMulOnScalar(b[n1:], x), vectorMulOnScalar(b[:n1], xinv))

	subProof, err := innerArgProof(&InnerArgumentPublic{
		N: n1,
		G: g1,
		H: h1,
		U: public.U,
	}, a1, b1)

	if err != nil {
		return nil, err
	}

	proof.A = subProof.A
	proof.B = subProof.B
	proof.L = append(proof.L, subProof.L...)
	proof.R = append(proof.R, subProof.R...)
	return proof, nil
}
