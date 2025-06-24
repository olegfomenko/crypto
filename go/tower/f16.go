package tower

// F16 represents the F8 field extension F16 = F8[x”]/(x”^2 + x'*x” + 1)
// x' \in F4[x']/(x'^2 + x' + 1)
// Representation costs 8 bit, value range [0,128)
type F16 struct {
	X1, X2 *F8
}

func (f *F16) Add(g *F16) *F16 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F16) Neg() *F16 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F8_10 = &F8{F4One(), F4Zero()}

func (f *F16) Mul(g *F16) *F16 {
	a0, a1 := NewF8().Set(f.X1), NewF8().Set(f.X2)
	b0, b1 := NewF8().Set(g.X1), NewF8().Set(g.X2)

	x1 := NewF8().Set(a0).Mul(b1).Add(
		NewF8().Set(a1).Mul(b0),
	).Add(
		NewF8().Set(a0).Mul(b0).Mul(F8_10),
	)

	x2 := NewF8().Set(a0).Mul(b0).Add(
		NewF8().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F16) Inv() *F16 {
	if f.Equal(F16Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF8().Set(f.X1)
	b := NewF8().Set(f.X2)

	d := NewF8().Set(a).Mul(b).Mul(F8_10).Add(
		NewF8().Set(b).Mul(b),
	).Add(
		NewF8().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F8_10)).Mul(d)
	return f
}

func (f *F16) Set(g *F16) *F16 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F16Zero() *F16 {
	return &F16{F8Zero(), F8Zero()}
}

func F16One() *F16 {
	return &F16{F8Zero(), F8One()}
}

func NewF16() *F16 {
	return F16One()
}

func (f *F16) Equal(g *F16) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF16() *F16 {
	return &F16{
		X1: RandomF8(),
		X2: RandomF8(),
	}
}
