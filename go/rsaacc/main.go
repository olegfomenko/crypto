package rsaacc

import (
	"crypto/rand"
	"github.com/olegfomenko/crypto/go/math"
	"math/big"
)

const KeySize = 128

const Exp = 65537

func mustGetBigPrime(bytes int) *big.Int {
	res, err := rand.Prime(rand.Reader, bytes*8)
	if err != nil {
		panic(err)
	}

	return res
}

func Gen() *big.Int {
	p, q := mustGetBigPrime(KeySize/2), mustGetBigPrime(KeySize/2)
	return p.Mul(p, q)
}

func Build(n, g *big.Int, list ...*big.Int) *big.Int {
	if len(list) == 0 {
		panic("Can not build accumulator for empty list")
	}

	for _, val := range list {
		ok, err := math.TestPrime(val)
		if err != nil {
			panic(err)
		}

		if !ok {
			panic("Accumulated values has to be prime")
		}
	}

	prod := big.NewInt(1)
	for _, val := range list {
		prod.Mul(prod, val)
	}

	return new(big.Int).Exp(g, prod, n)
}

func Prove(n, g *big.Int, pos int, list ...*big.Int) *big.Int {
	prod := big.NewInt(1)
	for i, val := range list {
		if i != pos {
			prod.Mul(prod, val)
		}
	}

	return new(big.Int).Exp(g, prod, n)
}

func Verify(n, witness, value, commit *big.Int) bool {
	return new(big.Int).Exp(witness, value, n).Cmp(commit) == 0
}
