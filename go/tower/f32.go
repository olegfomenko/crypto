package tower

// F32 represents the F16 field extension F32 = F16[x”]/(x”^2 + x'*x” + 1)
// x' \in F8[x']/(x'^2 + x' + 1)
// Representation costs 8 bit, value range [0,128)
type F32 struct {
	X1, X2 *F16
}

func (f *F32) Add(g *F32) *F32 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F32) Neg() *F32 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F16_10 = &F16{F8One(), F8Zero()}

func (f *F32) Mul(g *F32) *F32 {
	a0, a1 := NewF16().Set(f.X1), NewF16().Set(f.X2)
	b0, b1 := NewF16().Set(g.X1), NewF16().Set(g.X2)

	x1 := NewF16().Set(a0).Mul(b1).Add(
		NewF16().Set(a1).Mul(b0),
	).Add(
		NewF16().Set(a0).Mul(b0).Mul(F16_10),
	)

	x2 := NewF16().Set(a0).Mul(b0).Add(
		NewF16().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F32) Inv() *F32 {
	if f.Equal(F32Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF16().Set(f.X1)
	b := NewF16().Set(f.X2)

	d := NewF16().Set(a).Mul(b).Mul(F16_10).Add(
		NewF16().Set(b).Mul(b),
	).Add(
		NewF16().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F16_10)).Mul(d)
	return f
}

func (f *F32) Set(g *F32) *F32 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F32Zero() *F32 {
	return &F32{F16Zero(), F16Zero()}
}

func F32One() *F32 {
	return &F32{F16Zero(), F16One()}
}

func NewF32() *F32 {
	return F32One()
}

func (f *F32) Equal(g *F32) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF32() *F32 {
	return &F32{
		X1: RandomF16(),
		X2: RandomF16(),
	}
}
