package bp

import (
	"math/big"

	"github.com/cloudflare/bn256"
)

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
