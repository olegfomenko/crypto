package ec

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestECAdd(t *testing.T) {
	curveEth := secp256k1.S256()
	curveOleg := SECP256K1()

	_, x1, y1, err := elliptic.GenerateKey(curveEth, rand.Reader)
	if err != nil {
		panic(err)
	}

	_, x2, y2, err := elliptic.GenerateKey(curveEth, rand.Reader)
	if err != nil {
		panic(err)
	}

	xres1, yres1 := curveEth.Add(x1, y1, x2, y2)
	xres2, yres2 := curveOleg.Add(x1, y1, x2, y2)

	if xres1.Cmp(xres2) != 0 {
		panic("x result is not equal")
	}

	if yres1.Cmp(yres2) != 0 {
		panic("y result is not equal")
	}
}

func TestECDouble(t *testing.T) {
	curveEth := secp256k1.S256()
	curveOleg := SECP256K1()

	_, x1, y1, err := elliptic.GenerateKey(curveEth, rand.Reader)
	if err != nil {
		panic(err)
	}

	xres1, yres1 := curveEth.Double(x1, y1)
	xres2, yres2 := curveOleg.Double(x1, y1)

	if xres1.Cmp(xres2) != 0 {
		panic("x result is not equal")
	}

	if yres1.Cmp(yres2) != 0 {
		panic("y result is not equal")
	}
}

func TestECMul(t *testing.T) {
	curveEth := secp256k1.S256()
	curveOleg := SECP256K1()

	_, x1, y1, err := elliptic.GenerateKey(curveEth, rand.Reader)
	if err != nil {
		panic(err)
	}

	k := new(big.Int).SetInt64(1234567)

	xres1, yres1 := curveEth.ScalarMult(x1, y1, k.Bytes())
	xres2, yres2 := curveOleg.ScalarMult(x1, y1, k.Bytes())

	if xres1.Cmp(xres2) != 0 {
		panic("x result is not equal")
	}

	if yres1.Cmp(yres2) != 0 {
		panic("y result is not equal")
	}
}
