package tower

// F64 represents the F32 field extension F64 = F32[x”]/(x”^2 + x'*x” + 1)
// x' \in F16[x']/(x'^2 + x' + 1)
// Representation costs 8 bit, value range [0,128)
type F64 struct {
	X1, X2 *F32
}

func (f *F64) Add(g *F64) *F64 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F64) Neg() *F64 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F32_10 = &F32{F16One(), F16Zero()}

func (f *F64) Mul(g *F64) *F64 {
	a0, a1 := NewF32().Set(f.X1), NewF32().Set(f.X2)
	b0, b1 := NewF32().Set(g.X1), NewF32().Set(g.X2)

	x1 := NewF32().Set(a0).Mul(b1).Add(
		NewF32().Set(a1).Mul(b0),
	).Add(
		NewF32().Set(a0).Mul(b0).Mul(F32_10),
	)

	x2 := NewF32().Set(a0).Mul(b0).Add(
		NewF32().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F64) Inv() *F64 {
	if f.Equal(F64Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF32().Set(f.X1)
	b := NewF32().Set(f.X2)

	d := NewF32().Set(a).Mul(b).Mul(F32_10).Add(
		NewF32().Set(b).Mul(b),
	).Add(
		NewF32().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F32_10)).Mul(d)
	return f
}

func (f *F64) Set(g *F64) *F64 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F64Zero() *F64 {
	return &F64{F32Zero(), F32Zero()}
}

func F64One() *F64 {
	return &F64{F32Zero(), F32One()}
}

func NewF64() *F64 {
	return F64One()
}

func (f *F64) Equal(g *F64) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF64() *F64 {
	return &F64{
		X1: RandomF32(),
		X2: RandomF32(),
	}
}
