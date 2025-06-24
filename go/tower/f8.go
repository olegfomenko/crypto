package tower

// F8 represents the F4 field extension F8 = F4[x”]/(x”^2 + x'*x” + 1)
// x' \in F2[x']/(x'^2 + x' + 1)
// Representation costs 8 bit, value range [0,128)
type F8 struct {
	X1, X2 *F4
}

func (f *F8) Bytes() []byte {
	return []byte{
		f.X1.X1.X1.X<<7 + f.X1.X1.X2.X<<6 + f.X1.X2.X1.X<<5 + f.X1.X2.X2.X<<4 + f.X2.X1.X1.X<<3 + f.X2.X1.X2.X<<2 + f.X2.X2.X1.X<<1 + f.X2.X2.X2.X,
	}
}

func (f *F8) Add(g *F8) *F8 {
	f.X1.Add(g.X1)
	f.X2.Add(g.X2)
	return f
}

func (f *F8) Neg() *F8 {
	f.X1.Neg()
	f.X2.Neg()
	return f
}

var F4_10 = &F4{F2One(), F2Zero()}

func (f *F8) Mul(g *F8) *F8 {
	a0, a1 := NewF4().Set(f.X1), NewF4().Set(f.X2)
	b0, b1 := NewF4().Set(g.X1), NewF4().Set(g.X2)

	x1 := NewF4().Set(a0).Mul(b1).Add(
		NewF4().Set(a1).Mul(b0),
	).Add(
		NewF4().Set(a0).Mul(b0).Mul(F4_10),
	)

	x2 := NewF4().Set(a0).Mul(b0).Add(
		NewF4().Set(a1).Mul(b1),
	)

	f.X1.Set(x1)
	f.X2.Set(x2)
	return f
}

func (f *F8) Inv() *F8 {
	if f.Equal(F8Zero()) {
		panic("Multiplicative inverse does not exist for 0 in F_4")
	}

	a := NewF4().Set(f.X1)
	b := NewF4().Set(f.X2)

	d := NewF4().Set(a).Mul(b).Mul(F4_10).Add(
		NewF4().Set(b).Mul(b),
	).Add(
		NewF4().Set(a).Mul(a),
	).Inv()

	f.X1.Mul(d)
	f.X2.Add(a.Mul(F4_10)).Mul(d)
	return f
}

func (f *F8) Set(g *F8) *F8 {
	f.X1.Set(g.X1)
	f.X2.Set(g.X2)
	return f
}

func F8Zero() *F8 {
	return &F8{F4Zero(), F4Zero()}
}

func F8One() *F8 {
	return &F8{F4Zero(), F4One()}
}

func NewF8() *F8 {
	return F8One()
}

func (f *F8) Equal(g *F8) bool {
	return f.X1.Equal(g.X1) && f.X2.Equal(g.X2)
}

func RandomF8() *F8 {
	return &F8{
		X1: RandomF4(),
		X2: RandomF4(),
	}
}

func (f *F8) String() string {
	return f.X1.String() + f.X2.String()
}
