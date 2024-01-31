package bp

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

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
