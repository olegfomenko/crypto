package tower

// F128 represents the F64 field extension F128 = F64[x”]/(x”^2 + x'*x” + 1)
// x' \in F32[x']/(x'^2 + x' + 1)
// Representation costs 8 bit, value range [0,128)
type F128 struct {
	X1, X2 *F64
}

func (f *F128) Add(g *F128) *F128 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F128) Neg() *F128 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F64_10 = &F64{F32One(), F32Zero()}

func (f *F128) Mul(g *F128) *F128 {
	a0, a1 := NewF64().Set(f.X1), NewF64().Set(f.X2)
	b0, b1 := NewF64().Set(g.X1), NewF64().Set(g.X2)

	x1 := NewF64().Set(a0).Mul(b1).Add(
		NewF64().Set(a1).Mul(b0),
	).Add(
		NewF64().Set(a0).Mul(b0).Mul(F64_10),
	)

	x2 := NewF64().Set(a0).Mul(b0).Add(
		NewF64().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F128) Inv() *F128 {
	if f.Equal(F128Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF64().Set(f.X1)
	b := NewF64().Set(f.X2)

	d := NewF64().Set(a).Mul(b).Mul(F64_10).Add(
		NewF64().Set(b).Mul(b),
	).Add(
		NewF64().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F64_10)).Mul(d)
	return f
}

func (f *F128) Set(g *F128) *F128 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F128Zero() *F128 {
	return &F128{F64Zero(), F64Zero()}
}

func F128One() *F128 {
	return &F128{F64Zero(), F64One()}
}

func NewF128() *F128 {
	return F128One()
}

func (f *F128) Equal(g *F128) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF128() *F128 {
	return &F128{
		X1: RandomF64(),
		X2: RandomF64(),
	}
}
