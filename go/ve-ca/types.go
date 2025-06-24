package ve_ca

import (
	"crypto/rand"
	"github.com/cloudflare/bn256"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type F *big.Int
type G *bn256.G1

func MustRandomF() F {
	f, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}
	return f
}

func MustRandomG() G {
	_, g, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		panic(err)
	}
	return g
}

func MustSampleRandomF(k int) []F {
	res := make([]F, 0, k)
	for range k {
		res = append(res, MustRandomF())
	}

	return res
}

func GAdd(g1, g2 G) G {
	return new(bn256.G1).Add(g1, g2)
}

func GMul(s F, g G) G {
	return new(bn256.G1).ScalarMult(g, s)
}

func GBytes(g G) []byte {
	return (*bn256.G1)(g).Marshal()
}

func GArrBytes(g []G) []byte {
	var res []byte
	for i := range g {
		res = append(res, GBytes(g[i])...)
	}
	return res
}

func FSub(f1, f2 F) F {
	return new(big.Int).Mod(new(big.Int).Sub(f1, f2), bn256.Order)
}

func FMul(f1, f2 F) F {
	return new(big.Int).Mod(new(big.Int).Mul(f1, f2), bn256.Order)
}

func FPow(f1, f2 F) F {
	return new(big.Int).Exp(f1, f2, bn256.Order)
}

func FBytes(f F) []byte {
	return new(big.Int).Set(f).Bytes()
}

func FDiv(f1, f2 F) F {
	f2i := new(big.Int).ModInverse(f2, bn256.Order)
	return FMul(f1, f2i)
}

func FArrBytes(f []F) []byte {
	var res []byte
	for i := range f {
		res = append(res, FBytes(f[i])...)
	}
	return res
}

func FBit(f F, pos int) uint {
	return (*big.Int)(f).Bit(pos)
}

func Hash(b ...[]byte) F {
	h := crypto.Keccak256Hash(b...).Bytes()
	return new(big.Int).Mod(new(big.Int).SetBytes(h), bn256.Order)
}
