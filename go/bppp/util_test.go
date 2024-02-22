package bppp

import (
	"fmt"
	"math/big"
	"testing"
)

func TestPow(t *testing.T) {
	fmt.Println(pow(bint(2), -1))
	fmt.Println(inv(bint(2)))
}

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

func TestPolyVectorCalc(t *testing.T) {
	a := map[int][]*big.Int{
		1: ones(2),
		2: ones(2),
		3: ones(2),
	}

	// [1, 1] * 10 + [1, 1] * 100 + [1, 1] * 1000 = [1110, 1110]
	fmt.Println(polyVectorCalc(a, bint(10)))
}

func TestPolyCalc(t *testing.T) {
	a := map[int]*big.Int{
		1: bint(1),
		2: bint(1),
		3: bint(1),
	}

	fmt.Println(polyCalc(a, bint(10))) // 1110
}

func TestPolyVectorMul(t *testing.T) {
	a := map[int][]*big.Int{
		1: {bint(1), bint(2)},
		2: {bint(1), bint(3)},
	}

	b := map[int][]*big.Int{
		1: {bint(2), bint(1)},
		2: {bint(3), bint(1)},
	}

	// 1+1 = 4
	// 1 +2 = 5
	// 2 + 1 = 5
	// 2 + 2 = 6
	// res = [4, 10, 6]

	fmt.Println(polyVectorMul(a, b))
}

func TestPolyVectorAdd(t *testing.T) {
	a := map[int][]*big.Int{
		1: {bint(1), bint(2)},
		2: {bint(1), bint(3)},
	}

	b := map[int][]*big.Int{
		1: {bint(2), bint(1)},
		3: {bint(3), bint(1)},
	}

	// 1: 3, 3
	// 2: 1, 3
	// 3: 3, 1
	fmt.Println(polyVectorAdd(a, b))
}

func TestPolySub(t *testing.T) {
	a := map[int]*big.Int{
		1: bint(1),
		2: bint(4),
		4: bint(3),
	}

	b := map[int]*big.Int{
		1: bint(1),
		2: bint(2),
		3: bint(3),
	}

	// [0, 2, -3, 3]
	fmt.Println(polySub(a, b))
}

func TestPolyAdd(t *testing.T) {
	a := map[int]*big.Int{
		1: bint(1),
		2: bint(4),
		4: bint(3),
	}

	b := map[int]*big.Int{
		1: bint(1),
		2: bint(2),
		3: bint(3),
	}

	// [2, 6, 3, 3]
	fmt.Println(polyAdd(a, b))
}

func TestPolyMul(t *testing.T) {
	a := map[int]*big.Int{
		1: bint(1),
		2: bint(4),
	}

	b := map[int]*big.Int{
		1: bint(1),
		2: bint(2),
	}

	// 1 + 1 = 1
	// 1 + 2 = 2
	// 2 + 1 = 4
	// 2 + 2 = 8

	// [1, 6, 8]
	fmt.Println(polyMul(a, b))
}

func TestDiagInv(t *testing.T) {
	fmt.Println(diagInv(bint(1), 3))
	fmt.Println(diagInv(bint(1), 1))
	fmt.Println(diagInv(bint(10), 1))
	fmt.Println(diagInv(bint(10), 3))
}

func TestVectorMulOnMatrix(t *testing.T) {
	a := []*big.Int{bint(10), bint(5)}
	m := [][]*big.Int{
		{bint(1), bint(2)},
		{bint(3), bint(4)},
	}

	// a * m = [10*1 + 5*3, 10*2 + 5*4] = [25, 40]
	fmt.Println(vectorMulOnMatrix(a, m))
}
