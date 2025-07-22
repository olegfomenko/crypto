package gc

type GateType byte

const (
	GateXor GateType = iota
	GateAnd
	GateOr
)

type AbstractGate struct {
	X, Y     *AbstractGate
	Out      *AbstractGate
	Type     GateType
	Instance Gate
}

func (g *AbstractGate) IsXInput(input *AbstractGate) bool {
	return g.X == input
}

func (g *AbstractGate) IsYInput(input *AbstractGate) bool {
	return g.Y == input
}

func (g *AbstractGate) Initialized() bool {
	return g.Instance != nil
}

type Circuit struct {
	Inputs  []*AbstractGate
	Outputs []*AbstractGate
	results map[*AbstractGate]Label
}

func (c *Circuit) Run(inputs []Label) []Label {
	c.results = make(map[*AbstractGate]Label)

	for i, g := range c.Inputs {
		c.results[g] = g.Instance.Next(inputs[i*2], inputs[i*2+1])
		if g.Out != nil {
			c.run(g.Out)
		}
	}

	result := make([]Label, 0, len(c.Outputs))
	for _, g := range c.Outputs {
		result = append(result, c.results[g])
	}
	return result
}

func (c *Circuit) run(g *AbstractGate) {
	x, okX := c.results[g.X]
	y, okY := c.results[g.Y]

	if okX && okY {
		out := g.Instance.Next(x, y)
		c.results[g] = out
		if g.Out != nil {
			c.run(g.Out)
		}
	}

	return
}

func (c *Circuit) Build(out []Label) {
	for i, g := range c.Outputs {
		buildOutput(g, out[i*2], out[i*2+1])
	}
}

func buildOutput(g *AbstractGate, out0 Label, out1 Label) {
	if g.Initialized() {
		return
	}

	switch g.Type {
	case GateXor:
		g.Instance = NewXorGateWithOutput(out0, out1)
	case GateAnd:
		g.Instance = NewAndGateWithOut(out0, out1)
	case GateOr:
		g.Instance = NewOrGateWithOut(out0, out1)
	default:
		panic("invalid gate type")
	}

	if g.X != nil {
		build(g.X, g)
	}

	if g.Y != nil {
		build(g.Y, g)
	}
}

func build(g, parent *AbstractGate) {
	g.Out = parent

	if g.Initialized() {
		return
	}

	if g.Out.Initialized() {
		var out0, out1 Label

		if g.Out.IsXInput(g) {
			out0, out1 = getXInput(g.Out.Instance)
		} else {
			out0, out1 = getYInput(g.Out.Instance)
		}

		switch g.Type {
		case GateXor:
			g.Instance = NewXorGateWithOutput(out0, out1)
		case GateAnd:
			g.Instance = NewAndGateWithOut(out0, out1)
		case GateOr:
			g.Instance = NewOrGateWithOut(out0, out1)
		default:
			panic("invalid gate type")
		}
	}

	if g.X != nil {
		build(g.X, g)
	}

	if g.Y != nil {
		build(g.Y, g)
	}
}

func getXInput(g Gate) (Label, Label) {
	switch t := g.(type) {
	case *XorGate:
		return t.X0.L, t.X1.L
	case *GateTable:
		return t.X0.L, t.X1.L
	}

	panic("invalid gate type")
}

func getYInput(g Gate) (Label, Label) {
	switch t := g.(type) {
	case *XorGate:
		return t.Y0.L, t.Y1.L
	case *GateTable:
		return t.Y0.L, t.Y1.L
	}

	panic("invalid gate type")
}
