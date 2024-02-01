package bp

import (
	"crypto/rand"
	"math/big"

	"github.com/cloudflare/bn256"
)

func zeros(n int) []*big.Int {
	res := make([]*big.Int, n)
	for i := range res {
		res[i] = big.NewInt(0)
	}
	return res
}

func points(n int) []*bn256.G1 {
	res := make([]*bn256.G1, n)
	var err error
	for i := range res {
		if _, res[i], err = bn256.RandomG1(rand.Reader); err != nil {
			panic(err)
		}
	}
	return res
}

func vectorAdd(a []*big.Int, b []*big.Int) []*big.Int {
	if len(b) != len(a) {
		panic("invalid length")
	}

	res := make([]*big.Int, len(a))
	for i := 0; i < len(res); i++ {
		res[i] = add(a[i], b[i])
	}

	return res
}

func vectorMulOnScalar(a []*big.Int, c *big.Int) []*big.Int {
	res := make([]*big.Int, len(a))
	for i := range res {
		res[i] = mul(a[i], c)
	}
	return res
}

func vectorMul(a []*big.Int, b []*big.Int) *big.Int {
	if len(b) != len(a) {
		panic("invalid length")
	}

	res := big.NewInt(0)
	for i := 0; i < len(a); i++ {
		res = add(res, mul(a[i], b[i]))
	}
	return res
}

func vectorPointScalarMul(g []*bn256.G1, a []*big.Int) *bn256.G1 {
	if len(g) != len(a) {
		panic("invalid length for scalar mul")
	}

	res := new(bn256.G1).ScalarMult(g[0], a[0])
	for i := 1; i < len(g); i++ {
		res.Add(res, new(bn256.G1).ScalarMult(g[i], a[i]))
	}
	return res
}

func vectorPointMulOnScalar(g []*bn256.G1, a *big.Int) []*bn256.G1 {
	res := make([]*bn256.G1, len(g))
	for i := range res {
		res[i] = new(bn256.G1).ScalarMult(g[i], a)
	}
	return res
}

func hadamardMul(a, b []*bn256.G1) []*bn256.G1 {
	if len(a) != len(b) {
		panic("invalid length for scalar mul")
	}

	res := make([]*bn256.G1, len(a))
	for i := range res {
		res[i] = new(bn256.G1).Add(a[i], b[i])
	}
	return res
}

func add(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(x, y), bn256.Order)
}

func mul(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, y), bn256.Order)
}
