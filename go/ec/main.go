package ec

import (
	"crypto/elliptic"
	"math/big"
)

// Curve implements elliptic curve with equation: y^2 =x^3 + ax + b
type Curve struct {
	P       *big.Int
	N       *big.Int
	A       *big.Int
	B       *big.Int
	Gx, Gy  *big.Int
	BitSize int
}

// SECP256K1 returns an Ethereum secp256k1 curve
func SECP256K1() *Curve {
	curve := &Curve{}

	curve.P, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 0)
	curve.N, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)

	curve.A = new(big.Int).SetInt64(0)
	curve.B = new(big.Int).SetInt64(7)

	curve.Gx, _ = new(big.Int).SetString("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 0)
	curve.Gy, _ = new(big.Int).SetString("0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 0)
	curve.BitSize = 256
	return curve
}

var _ elliptic.Curve = &Curve{}

func (c *Curve) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P:       c.P,
		N:       c.N,
		B:       c.B,
		Gx:      c.Gx,
		Gy:      c.Gy,
		BitSize: c.BitSize,
		Name:    "secp256k1",
	}
}

func (c *Curve) IsOnCurve(x, y *big.Int) bool {
	if x == nil && y == nil {
		return true
	}

	yy := mul(y, y, c.P)
	xxx := mul(mul(x, x, c.P), x, c.P)
	ax := mul(c.A, x, c.P)
	return yy.Cmp(add(xxx, ax, c.P)) == 0
}

func (c *Curve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	if x1 == nil && y1 == nil {
		return x2, y2
	}

	if x2 == nil && y2 == nil {
		return x1, y1
	}

	s := div(sub(y1, y2, c.P), sub(x1, x2, c.P), c.P)
	x = sub(sub(mul(s, s, c.P), x1, c.P), x2, c.P)
	y = sub(mul(s, sub(x1, x, c.P), c.P), y1, c.P)
	return
}

func (c *Curve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	if x1 == nil && y1 == nil {
		return nil, nil
	}

	s := div(add(mul(fromInt(3), mul(x1, x1, c.P), c.P), c.A, c.P), mul(fromInt(2), y1, c.P), c.P)
	x = sub(mul(s, s, c.P), mul(fromInt(2), x1, c.P), c.P)
	y = sub(mul(s, sub(x1, x, c.P), c.P), y1, c.P)
	return
}

func (c *Curve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	if len(k) > 32 {
		panic("K have to be 256 bits maximum")
	}

	if x1 == nil && y1 == nil {
		return nil, nil
	}

	scalar := new(big.Int).Mod(new(big.Int).SetBytes(k), c.N)
	bits := scalar.Text(2)

	for i := len(bits) - 1; i >= 0; i-- {
		bit := bits[i]
		if bit == '1' {
			x, y = c.Add(x, y, x1, y1)
		}

		x1, y1 = c.Double(x1, y1)
	}

	return
}

func (c *Curve) mult(k, x1, y1 *big.Int) (x, y *big.Int) {
	if k.Cmp(fromInt(0)) == 0 {
		return nil, nil
	}

	if k.Cmp(fromInt(1)) == 0 {
		return x1, y1
	}

	x, y = c.mult(new(big.Int).Div(k, fromInt(2)), x1, y1)
	x, y = c.Double(x, y)

	if k.Bytes()[0] == 0 {
		return x, y
	}

	return c.Add(x, y, x1, y1)
}

func (c *Curve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return c.ScalarMult(c.Gx, c.Gx, k)
}

func add(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(x, y), mod)
}

func sub(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Sub(x, y), mod)
}

func mul(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, y), mod)
}

func div(x *big.Int, y *big.Int, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, new(big.Int).ModInverse(y, mod)), mod)
}

func fromInt(val int64) *big.Int {
	return new(big.Int).SetInt64(val)
}
