package bppp

import (
	"fmt"
	"math/big"
	"testing"
)

func TestTensorMul(t *testing.T) {
	a := powvector(bint(10), 2)
	b := powvector(bint(100), 3)
	// 1 10
	// 1 100 10000

	// 1 10 100 1000 10000 100000
	fmt.Println(vectorTensorMul(a, b))
}

func TestPolyVectorMulWeight(t *testing.T) {
	mu := bint(10)
	a := map[int][]*big.Int{
		1: values(2),
		2: values(2),
		3: values(2),
	}

	x := values(1)[0]

	ax := polyVectorCalc(a, x)
	res1 := weightVectorMul(ax, ax, mu)
	fmt.Println(res1)

	a2mu := polyVectorMulWeight(a, a, mu)
	res2 := polyCalc(a2mu, x)
	fmt.Println(res2)

	a2mu2 := polyVectorMulWeight2(a, a, mu)
	res3 := polyCalc(a2mu2, x)
	fmt.Println(res3)
}
