package gc

import "bytes"

type Gate interface {
	Next(x Label, y Label) Label
}

type GateLabel struct {
	Hash [32]byte
	L    Label
}

var _ Gate = &GateTable{}

type GateTable struct {
	X0, X1, Y0, Y1 GateLabel
	Out00          [32]byte
	Out01          [32]byte
	Out10          [32]byte
	Out11          [32]byte
}

func (g *GateTable) Next(x Label, y Label) Label {
	xHash := Hash(x)
	yHash := Hash(y)

	if bytes.Equal(xHash[:], g.X0.Hash[:]) {
		if bytes.Equal(yHash[:], g.Y0.Hash[:]) {
			return Encrypt(x, y, g.Out00)
		} else {
			return Encrypt(x, y, g.Out01)
		}
	} else {
		if bytes.Equal(yHash[:], g.Y0.Hash[:]) {
			return Encrypt(x, y, g.Out10)
		} else {
			return Encrypt(x, y, g.Out11)
		}
	}

	panic("invalid gate oppenings")
}

func NewAndGate(x0, x1, y0, y1 Label, out0, out1 Label) Gate {
	return &GateTable{
		X0:    GateLabel{Hash(x0), x0},
		X1:    GateLabel{Hash(x1), x1},
		Y0:    GateLabel{Hash(y0), y0},
		Y1:    GateLabel{Hash(y1), y1},
		Out00: Encrypt(x0, y0, out0),
		Out01: Encrypt(x0, y1, out0),
		Out10: Encrypt(x1, y0, out0),
		Out11: Encrypt(x1, y1, out1),
	}
}

func NewAndGateWithOut(out0, out1 Label) Gate {
	x0 := RL()
	x1 := RL()
	y0 := RL()
	y1 := RL()
	return NewAndGate(x0, x1, y1, y0, out0, out1)
}

func NewOrGate(x0, x1, y0, y1 Label, out0, out1 Label) Gate {
	return &GateTable{
		X0:    GateLabel{Hash(x0), x0},
		X1:    GateLabel{Hash(x1), x1},
		Y0:    GateLabel{Hash(y0), y0},
		Y1:    GateLabel{Hash(y1), y1},
		Out00: Encrypt(x0, y0, out0),
		Out01: Encrypt(x0, y1, out1),
		Out10: Encrypt(x1, y0, out1),
		Out11: Encrypt(x1, y1, out1),
	}
}

func NewOrGateWithOut(out0, out1 Label) Gate {
	x0 := RL()
	x1 := RL()
	y0 := RL()
	y1 := RL()
	return NewOrGate(x0, x1, y1, y0, out0, out1)
}

var _ Gate = &XorGate{}

type XorGate struct {
	X0, X1, Y0, Y1 GateLabel
}

func (g *XorGate) Next(x Label, y Label) Label {
	return Xor(x, y)
}

func NewXorGate(x0, y0, r Label) Gate {
	x1 := Xor(x0, r)
	y1 := Xor(y0, r)

	return &XorGate{
		X0: GateLabel{Hash(x0), x0},
		X1: GateLabel{Hash(x1), x1},
		Y0: GateLabel{Hash(y0), y0},
		Y1: GateLabel{Hash(y1), y1},
	}
}

func NewXorGateWithOutput(out0, out1 Label) Gate {
	x0 := RL()
	y0 := Xor(out0, x0) // x0^y0 = out0 => y0 = out0 ^ x0

	r := Xor(out0, out1)
	// x1 := Xor(x0, r) // x1 = x0 ^ out0 ^ out1
	// y1 := Xor(y0, r) // y1 = out0 ^ x0 ^ out0 ^ out1 = x0 ^ out1

	// x0 ^ y0 = x0 ^ out0 ^ x0 = out0
	// x1 ^ y1 = x0 ^ out0 ^ out1 ^ x0 ^ out1 = out0
	// x1 ^ y0 = x0 ^ out0 ^ out1 ^ out0 ^ x0 = out1
	// x0 ^ y1 = x0 ^ x0 ^ out1 = out1

	return NewXorGate(x0, y0, r)
}
