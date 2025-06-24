package tower

import (
	"crypto/rand"
	"fmt"
)

// F1 represents the $F_2$ field structure
// Representation costs 1 bit, value range [0,2)
type F1 struct {
	X byte
}

func (f *F1) Add(g *F1) *F1 {
	f.X = f.X ^ g.X
	return f
}

func (f *F1) Neg() *F1 {
	return f
}

func (f *F1) Mul(g *F1) *F1 {
	f.X = f.X & g.X
	return f
}

func (f *F1) Inv() *F1 {
	if f.X == 0 {
		panic("Multiplicative inverse does not exist for 0 in F_2")
	}
	return f
}

func (f *F1) Set(g *F1) *F1 {
	f.X = g.X
	return f
}

func F1Zero() *F1 {
	return &F1{0}
}

func F1One() *F1 {
	return &F1{1}
}

func NewF1() *F1 {
	return F1One()
}

func (f *F1) Equal(g *F1) bool {
	return f.X == g.X
}

func RandomF1() *F1 {
	buf := make([]byte, 1)
	n, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	if n != 1 {
		panic("something went wrong")
	}

	return &F1{byte(buf[0] % 2)}
}

func (f *F1) String() string {
	return fmt.Sprint(f.X)
}
