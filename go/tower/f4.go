package tower

// F4 represents the F2 field extension F4 = F2[x”]/(x”^2 + x'*x” + 1)
// x' \in F1[x']/(x'^2 + x' + 1)
// Representation costs 4 bit, value range [0,16)
type F4 struct {
	X1, X2 *F2
}

func (f *F4) Add(g *F4) *F4 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F4) Neg() *F4 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F2_10 = &F2{F1One(), F1Zero()}

func (f *F4) Mul(g *F4) *F4 {
	a0, a1 := NewF2().Set(f.X1), NewF2().Set(f.X2)
	b0, b1 := NewF2().Set(g.X1), NewF2().Set(g.X2)

	x1 := NewF2().Set(a0).Mul(b1).Add(
		NewF2().Set(a1).Mul(b0),
	).Add(
		NewF2().Set(a0).Mul(b0).Mul(F2_10),
	)

	x2 := NewF2().Set(a0).Mul(b0).Add(
		NewF2().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F4) Inv() *F4 {
	if f.Equal(F4Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF2().Set(f.X1)
	b := NewF2().Set(f.X2)

	d := NewF2().Set(a).Mul(b).Mul(F2_10).Add(
		NewF2().Set(b).Mul(b),
	).Add(
		NewF2().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F2_10)).Mul(d)
	return f
}

func (f *F4) Set(g *F4) *F4 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F4Zero() *F4 {
	return &F4{F2Zero(), F2Zero()}
}

func F4One() *F4 {
	return &F4{F2Zero(), F2One()}
}

func NewF4() *F4 {
	return F4One()
}

func (f *F4) Equal(g *F4) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF4() *F4 {
	return &F4{
		X1: RandomF2(),
		X2: RandomF2(),
	}
}

func (f *F4) String() string {
	return f.X1.String() + f.X2.String()
}
