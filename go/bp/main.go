package bp

import (
	"math/big"

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

func productCom(g []*bn256.G1, h []*bn256.G1, u *bn256.G1, a, b []*big.Int) *bn256.G1 {
	p := vecCom(g, h, a, b)
	p.Add(p, new(bn256.G1).ScalarMult(u, vectorMul(a, b)))
	return p
}

func InnerProductH(n int, g []*bn256.G1, h []*bn256.G1, u *bn256.G1, a, a1, b, b1 []*big.Int, c *big.Int) *bn256.G1 {
	if n%2 != 0 {
		panic("invalid n")
	}

	n1 := n / 2
	res := vectorPointScalarMul(g[:n1], a)
	res.Add(res, vectorPointScalarMul(g[n1:], a1))
	res.Add(res, vectorPointScalarMul(h[:n1], b))
	res.Add(res, vectorPointScalarMul(h[n1:], b1))
	res.Add(res, new(bn256.G1).ScalarMult(u, c))
	return res
}
