package bp

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

func TestInnerArgumentProofRecursive(t *testing.T) {
	// public values
	const n = 4

	g := points(n)
	h := points(n)
	u := points(1)

	// private values
	a := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(10), big.NewInt(1)}
	b := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(10)}

	// public commitment P: p = g^a+h^b+u^<a, b>
	p := P(g, h, u[0], a, b)

	protocol2(n, g, h, u[0], p, a, b)
}

func protocol2(n int, g, h []*bn256.G1, u *bn256.G1, p *bn256.G1, a, b []*big.Int) {
	fmt.Println("Running protocol2() in n =", n)

	if n == 1 {
		// send a, b to verifier
		ga := new(bn256.G1).ScalarMult(g[0], a[0])
		hb := new(bn256.G1).ScalarMult(h[0], b[0])
		uc := new(bn256.G1).ScalarMult(u, mul(a[0], b[0]))

		p1 := new(bn256.G1).Add(ga, hb)
		p1.Add(p1, uc)

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

	L := vectorPointScalarMul(g[n1:], a[:n1])
	L.Add(L, vectorPointScalarMul(h[:n1], b[n1:]))
	L.Add(L, new(bn256.G1).ScalarMult(u, cl))

	R := vectorPointScalarMul(g[:n1], a[n1:])
	R.Add(R, vectorPointScalarMul(h[n1:], b[:n1]))
	R.Add(R, new(bn256.G1).ScalarMult(u, cr))

	// Send L, R to verifier

	// Verifier generates challenge and sends to prover
	x, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}

	xinv := new(big.Int).ModInverse(x, bn256.Order)
	x2 := mul(x, x)
	x2inv := new(big.Int).ModInverse(x2, bn256.Order)

	// Both verifier and prover computes new public values

	g1 := hadamardMul(vectorPointMulOnScalar(g[:n1], xinv), vectorPointMulOnScalar(g[n1:], x))
	h1 := hadamardMul(vectorPointMulOnScalar(h[:n1], x), vectorPointMulOnScalar(h[n1:], xinv))

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
