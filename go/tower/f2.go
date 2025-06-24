package tower

// F2 represents the F1 field extension F2 = F1[x]/(x^2 + x + 1)
// Each F2 element is exactly ax+b where a, b \in F1
// Here X1 is a and X2 is b
// Representation costs 2 bit, value range [0,4)
type F2 struct {
	X1, X2 *F1
}

func (f *F2) Add(g *F2) *F2 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F2) Neg() *F2 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

// $(a_1x+b_1)(a_2x+b_2) = a_1a_2x^2 + (a_1b_2 + a_2_b_1)x + b_1b_2$
// Let's represent it as $Ax^2 + Bx + C$, where $A = a_1a_2$, $B = a_1b_2 + a_2_b_1$, and $C = b_1b_2$
// Then, adter taking by modulo $x^2 + x + 1$ we receive $(B-A)x + (C-A)$
func (f *F2) Mul(g *F2) *F2 {
	a0, a1 := NewF1().Set(f.X1), NewF1().Set(f.X2)
	b0, b1 := NewF1().Set(g.X1), NewF1().Set(g.X2)

	x1 := NewF1().Set(a0).Mul(b1).Add(
		NewF1().Set(a1).Mul(b0),
	).Add(
		NewF1().Set(a0).Mul(b0),
	)

	x2 := NewF1().Set(a0).Mul(b0).Add(
		NewF1().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F2) Inv() *F2 {
	if f.Equal(F2Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF1().Set(f.X1)
	b := NewF1().Set(f.X2)

	d := NewF1().Set(a).Mul(b).Add(
		NewF1().Set(b).Mul(b),
	).Add(
		NewF1().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a).Mul(d)
	return f
}

func (f *F2) Set(g *F2) *F2 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F2Zero() *F2 {
	return &F2{F1Zero(), F1Zero()}
}

func F2One() *F2 {
	return &F2{F1Zero(), F1One()}
}

func NewF2() *F2 {
	return F2One()
}

func (f *F2) Equal(g *F2) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF2() *F2 {
	return &F2{
		X1: RandomF1(),
		X2: RandomF1(),
	}
}

func (f *F2) String() string {
	return f.X1.String() + f.X2.String()
}
