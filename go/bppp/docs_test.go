package bppp

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

func TestWNLA(t *testing.T) {
	// Public
	const N = 4

	g := points(1)[0]
	G := points(N)

	H := points(N)

	c := values(N)

	ro := values(1)[0]
	mu := mul(ro, ro)

	// Private
	l := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(10), big.NewInt(1)}
	n := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(10)}

	// Com
	v := add(vectorMul(c, l), wightVectorMul(n, n, mu))
	C := new(bn256.G1).ScalarMult(g, v)
	C.Add(C, vectorPointScalarMul(H, l))
	C.Add(C, vectorPointScalarMul(G, n))

	wnla(g, G, H, c, C, ro, mu, l, n)
}

func wnla(g *bn256.G1, G, H []*bn256.G1, c []*big.Int, C *bn256.G1, ro, mu *big.Int, l, n []*big.Int) {
	roinv := new(big.Int).ModInverse(ro, bn256.Order)

	if len(l)+len(n) < 6 {

		// Prover sends l, n to Verifier
		// Next verifier computes:
		_v := add(vectorMul(c, l), wightVectorMul(n, n, mu))

		_C := new(bn256.G1).ScalarMult(g, _v)
		_C.Add(_C, vectorPointScalarMul(H, l))
		_C.Add(_C, vectorPointScalarMul(G, n))

		if !bytes.Equal(_C.Marshal(), C.Marshal()) {
			panic("Failed to verify!")
		}

		fmt.Println("Verified!")
		return
	}

	// Verifier selects random challenge
	y := values(1)[0]

	// Prover calculates new reduced values, vx and vr and sends X, R to verifier
	c0, c1 := reduceVector(c)
	l0, l1 := reduceVector(l)
	n0, n1 := reduceVector(n)
	G0, G1 := reducePoints(G)
	H0, H1 := reducePoints(H)

	l_ := vectorAdd(l0, vectorMulOnScalar(l1, y))
	n_ := vectorAdd(vectorMulOnScalar(n0, roinv), vectorMulOnScalar(n1, y))

	//v_ := add(vectorMul(c_, l_), wightVectorMul(n_, n_, mul(mu, mu)))

	vx := add(
		mul(wightVectorMul(n0, n1, mul(mu, mu)), mul(big.NewInt(2), roinv)),
		add(vectorMul(c0, l1), vectorMul(c1, l0)),
	)

	vr := add(wightVectorMul(n1, n1, mul(mu, mu)), vectorMul(c1, l1))

	X := new(bn256.G1).ScalarMult(g, vx)
	X.Add(X, vectorPointScalarMul(H0, l1))
	X.Add(X, vectorPointScalarMul(H1, l0))
	X.Add(X, vectorPointScalarMul(G0, vectorMulOnScalar(n1, ro)))
	X.Add(X, vectorPointScalarMul(G1, vectorMulOnScalar(n0, roinv)))

	R := new(bn256.G1).ScalarMult(g, vr)
	R.Add(R, vectorPointScalarMul(H1, l1))
	R.Add(R, vectorPointScalarMul(G1, n1))

	// Submit R, X to Verifier

	// Both computes
	H_ := vectorPointsAdd(H0, vectorPointMulOnScalar(H1, y))
	G_ := vectorPointsAdd(vectorPointMulOnScalar(G0, ro), vectorPointMulOnScalar(G1, y))
	c_ := vectorAdd(c0, vectorMulOnScalar(c1, y))

	ro_ := mu
	mu_ := mul(mu, mu)

	C_ := new(bn256.G1).Set(C)
	C_.Add(C_, new(bn256.G1).ScalarMult(X, y))
	C_.Add(C_, new(bn256.G1).ScalarMult(R, sub(mul(y, y), big.NewInt(1))))

	// Recursive run
	wnla(g, G_, H_, c_, C_, ro_, mu_, l_, n_)
}
