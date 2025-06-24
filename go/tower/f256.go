package tower

// F256 represents the F128 field extension F256 = F128[x”]/(x”^2 + x'*x” + 1)
// x' \in F64[x']/(x'^2 + x' + 1)
// Representation costs 8 bit, value range [0,128)
type F256 struct {
	X1, X2 *F128
}

func (f *F256) Add(g *F256) *F256 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F256) Neg() *F256 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F128_10 = &F128{F64One(), F64Zero()}

func (f *F256) Mul(g *F256) *F256 {
	a0, a1 := NewF128().Set(f.X1), NewF128().Set(f.X2)
	b0, b1 := NewF128().Set(g.X1), NewF128().Set(g.X2)

	x1 := NewF128().Set(a0).Mul(b1).Add(
		NewF128().Set(a1).Mul(b0),
	).Add(
		NewF128().Set(a0).Mul(b0).Mul(F128_10),
	)

	x2 := NewF128().Set(a0).Mul(b0).Add(
		NewF128().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F256) Inv() *F256 {
	if f.Equal(F256Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF128().Set(f.X1)
	b := NewF128().Set(f.X2)

	d := NewF128().Set(a).Mul(b).Mul(F128_10).Add(
		NewF128().Set(b).Mul(b),
	).Add(
		NewF128().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F128_10)).Mul(d)
	return f
}

func (f *F256) Set(g *F256) *F256 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F256Zero() *F256 {
	return &F256{F128Zero(), F128Zero()}
}

func F256One() *F256 {
	return &F256{F128Zero(), F128One()}
}

func NewF256() *F256 {
	return F256One()
}

func (f *F256) Equal(g *F256) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF256() *F256 {
	return &F256{
		X1: RandomF128(),
		X2: RandomF128(),
	}
}
