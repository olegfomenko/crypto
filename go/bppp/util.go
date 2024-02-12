package bppp

import (
	"crypto/rand"
	"math/big"

	"github.com/cloudflare/bn256"
)

func values(n int) []*big.Int {
	res := make([]*big.Int, n)
	var err error

	for i := range res {
		res[i], err = rand.Int(rand.Reader, bn256.Order)
		if err != nil {
			panic(err)
		}
	}

	return res
}

func zeros(n int) []*big.Int {
	res := make([]*big.Int, n)
	for i := range res {
		res[i] = big.NewInt(0)
	}
	return res
}

func powvector(v *big.Int, a int) []*big.Int {
	val := big.NewInt(1)
	res := make([]*big.Int, a)
	for i := range res {
		res[i] = val
		val = mul(val, v)
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

func inv(x *big.Int) *big.Int {
	return new(big.Int).ModInverse(x, bn256.Order)
}

func zeroIfNil(x *big.Int) *big.Int {
	if x == nil {
		return bint(0)
	}
	return x
}

func add(x *big.Int, y *big.Int) *big.Int {
	x = zeroIfNil(x)
	y = zeroIfNil(y)
	return new(big.Int).Mod(new(big.Int).Add(x, y), bn256.Order)
}

func sub(x *big.Int, y *big.Int) *big.Int {
	x = zeroIfNil(x)
	y = zeroIfNil(y)
	return new(big.Int).Mod(new(big.Int).Sub(x, y), bn256.Order)
}

func mul(x *big.Int, y *big.Int) *big.Int {
	if x == nil || y == nil {
		return bint(0)
	}

	return new(big.Int).Mod(new(big.Int).Mul(x, y), bn256.Order)
}

func pow(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Exp(x, y, bn256.Order)
}

func bint(v int) *big.Int {
	return new(big.Int).Mod(new(big.Int).SetInt64(int64(v)), bn256.Order)
}

func frac(a, b int) *big.Int {
	return mul(bint(a), inv(bint(b)))
}

func bbool(v bool) *big.Int {
	if v {
		return new(big.Int).SetInt64(1)
	}

	return new(big.Int).SetInt64(0)
}

func vectorTensorMul(a, b []*big.Int) []*big.Int {
	res := make([]*big.Int, 0, len(a)*len(b))

	for i := range a {
		res = append(res, vectorMulOnScalar(b, a[i])...)
	}
	return res
}

func vectorMul(a []*big.Int, b []*big.Int) *big.Int {
	for len(a) < len(b) {
		a = append(a, bint(0))
	}

	for len(b) < len(a) {
		b = append(b, bint(0))
	}

	res := big.NewInt(0)
	for i := 0; i < len(a); i++ {
		res = add(res, mul(a[i], b[i]))
	}
	return res
}

func weightVectorMul(a []*big.Int, b []*big.Int, mu *big.Int) *big.Int {
	for len(a) < len(b) {
		a = append(a, bint(0))
	}

	for len(b) < len(a) {
		b = append(b, bint(0))
	}

	res := big.NewInt(0)
	exp := new(big.Int).Set(mu)

	for i := 0; i < len(a); i++ {
		res = add(res, mul(mul(a[i], b[i]), exp))
		exp = mul(exp, mu)
	}
	return res
}

func vectorPointScalarMul(g []*bn256.G1, a []*big.Int) *bn256.G1 {
	if len(g) == 0 {
		return new(bn256.G1).ScalarBaseMult(bint(0))
	}

	for len(a) < len(g) {
		a = append(a, bint(0))
	}

	res := new(bn256.G1).ScalarMult(g[0], a[0])
	for i := 1; i < len(g); i++ {
		res.Add(res, new(bn256.G1).ScalarMult(g[i], a[i]))
	}
	return res
}

func reduceVector(v []*big.Int) ([]*big.Int, []*big.Int) {
	res0 := make([]*big.Int, 0, len(v)/2)
	res1 := make([]*big.Int, 0, len(v)/2)

	for i := range v {
		if i%2 == 0 {
			res0 = append(res0, v[i])
		} else {
			res1 = append(res1, v[i])
		}
	}

	return res0, res1
}

func reducePoints(v []*bn256.G1) ([]*bn256.G1, []*bn256.G1) {
	res0 := make([]*bn256.G1, 0, len(v)/2)
	res1 := make([]*bn256.G1, 0, len(v)/2)

	for i := range v {
		if i%2 == 0 {
			res0 = append(res0, v[i])
		} else {
			res1 = append(res1, v[i])
		}
	}

	return res0, res1
}

func vectorPointMulOnScalar(g []*bn256.G1, a *big.Int) []*bn256.G1 {
	res := make([]*bn256.G1, len(g))
	for i := range res {
		res[i] = new(bn256.G1).ScalarMult(g[i], a)
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

func vectorAdd(a []*big.Int, b []*big.Int) []*big.Int {
	for len(a) < len(b) {
		a = append(a, bint(0))
	}

	for len(b) < len(a) {
		b = append(b, bint(0))
	}

	res := make([]*big.Int, len(a))
	for i := 0; i < len(res); i++ {
		res[i] = add(a[i], b[i])
	}

	return res
}

func vectorSub(a []*big.Int, b []*big.Int) []*big.Int {
	for len(a) < len(b) {
		a = append(a, bint(0))
	}

	for len(b) < len(a) {
		b = append(b, bint(0))
	}

	res := make([]*big.Int, len(a))
	for i := 0; i < len(res); i++ {
		res[i] = sub(a[i], b[i])
	}

	return res
}

func vectorPointsAdd(a, b []*bn256.G1) []*bn256.G1 {
	for len(a) < len(b) {
		a = append(a, new(bn256.G1).ScalarBaseMult(bint(0)))
	}

	for len(b) < len(a) {
		b = append(b, new(bn256.G1).ScalarBaseMult(bint(0)))
	}

	res := make([]*bn256.G1, len(a))
	for i := range res {
		res[i] = new(bn256.G1).Add(a[i], b[i])
	}
	return res
}

func vectorMulOnMatrix(a []*big.Int, m [][]*big.Int) []*big.Int {
	var res []*big.Int

	for j := 0; j < len(m[0]); j++ {
		var column []*big.Int

		for i := 0; i < len(m); i++ {
			column = append(column, m[i][j])
		}

		res = append(res, vectorMul(a, column))
	}

	return res
}

func diagInv(x *big.Int, n int) [][]*big.Int {
	var res [][]*big.Int = make([][]*big.Int, n)
	val := big.NewInt(1)
	inv := new(big.Int).ModInverse(x, bn256.Order)

	for i := 0; i < n; i++ {
		res[i] = make([]*big.Int, n)

		for j := 0; j < n; j++ {
			res[i][j] = big.NewInt(0)

			if i == j {
				res[i][j] = val
				val = mul(val, inv)
			}
		}
	}

	return res
}

func polyMul(a, b map[int]*big.Int) map[int]*big.Int { // res dimension will be len(a) + len(b) - 1
	res := make(map[int]*big.Int)

	if len(a) == 0 || len(b) == 0 {
		panic("invalid polynom")
	}

	for i := range a {
		for j := range b {
			res[i+j] = add(res[i+j], mul(a[i], b[j]))
		}
	}

	return res
}

func polyAdd(a, b map[int]*big.Int) map[int]*big.Int { // res dimension will be max(len(a), len(b))
	res := make(map[int]*big.Int)

	for k, v := range a {
		res[k] = new(big.Int).Set(v)
	}

	for k, v := range b {
		res[k] = add(res[k], v)
	}

	return res
}

func polySub(a, b map[int]*big.Int) map[int]*big.Int { // res dimension will be max(len(a), len(b))
	res := make(map[int]*big.Int)

	for k, v := range a {
		res[k] = new(big.Int).Set(v)
	}

	for k, v := range b {
		res[k] = sub(res[k], v)
	}

	return res
}

func polyVectorAdd(a, b map[int][]*big.Int) map[int][]*big.Int { // res dimension will be max(len(a), len(b))
	res := make(map[int][]*big.Int)

	for k, v := range a {
		res[k] = v
	}

	for k, v := range b {
		res[k] = vectorAdd(res[k], v)
	}

	return res
}

func polyVectorMul(a, b map[int][]*big.Int) map[int]*big.Int { // res dimension will be len(a) + len(b) - 1
	res := make(map[int]*big.Int)

	if len(a) == 0 || len(b) == 0 {
		panic("invalid polynom")
	}

	for i := range a {
		for j := range b {
			res[i+j] = add(res[i+j], vectorMul(a[i], b[j]))
		}
	}

	return res
}

func polyCalc(poly map[int]*big.Int, x *big.Int) *big.Int {
	res := bint(0)
	for k, v := range poly {
		res = add(res, mul(v, pow(x, bint(k))))
	}
	return res
}

func polyVectorCalc(poly map[int][]*big.Int, x *big.Int) []*big.Int {
	var res []*big.Int
	for k, v := range poly {
		res = vectorAdd(res, vectorMulOnScalar(v, pow(x, bint(k))))
	}
	return res
}

func polyVectorMulWeight(a, b map[int][]*big.Int, mu *big.Int) map[int]*big.Int { // res dimension will be len(a) + len(b) - 1
	res := make(map[int]*big.Int)

	if len(a) == 0 || len(b) == 0 {
		panic("invalid polynom")
	}

	for i := range a {
		for j := range b {
			res[i+j] = add(weightVectorMul(a[i], b[j], mu), res[i+j])
		}
	}

	return res
}

// Only testing functions

func splitVectorPoly(a map[int][]*big.Int) []map[int]*big.Int {
	var res []map[int]*big.Int

	dim := 0
	for _, v := range a {
		dim = len(v)
		break
	}

	for i := 0; i < dim; i++ {
		res = append(res, make(map[int]*big.Int))
		for k, v := range a {
			res[i][k] = v[i]
		}
	}

	return res
}

func polyMulOnScalar(a map[int]*big.Int, mu *big.Int) map[int]*big.Int {
	res := make(map[int]*big.Int)
	for k, v := range a {
		res[k] = mul(v, mu)
	}
	return res
}

func polyVectorMulWeight2(a, b map[int][]*big.Int, mu *big.Int) map[int]*big.Int {
	res := make(map[int]*big.Int)

	polyA := splitVectorPoly(a)
	polyB := splitVectorPoly(b)

	if len(polyA) != len(polyB) {
		panic("can not mull")
	}

	exp := new(big.Int).Set(mu)

	for i := 0; i < len(polyA); i++ {
		res = polyAdd(res, polyMulOnScalar(polyMul(polyA[i], polyB[i]), exp))
		exp = mul(exp, mu)
	}

	return res
}
