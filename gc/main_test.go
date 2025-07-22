package gc

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGateGeneration(t *testing.T) {
	// 1 \
	//     AND
	// 2 /      \
	//            5
	//              XOR -> 7
	//            6
	// 3 \      /
	//     OR
	// 4 /

	out0 := Label{}
	out1 := Label{}
	out1[0] = 1

	gates := []Gate{
		NewXorGateWithOutput(out0, out1),
	}

	gates = append(gates, NewAndGateWithOut(gates[0].(*XorGate).X0.L, gates[0].(*XorGate).X1.L))
	gates = append(gates, NewOrGateWithOut(gates[0].(*XorGate).Y0.L, gates[0].(*XorGate).Y1.L))

	// Imagine input 0111
	l1 := gates[1].(*GateTable).X0
	l2 := gates[1].(*GateTable).Y1
	l3 := gates[2].(*GateTable).X1
	l4 := gates[2].(*GateTable).Y1

	l5 := gates[1].Next(l1.L, l2.L)
	l6 := gates[2].Next(l3.L, l4.L)

	result := gates[0].Next(l5, l6)

	fmt.Println(bytes.Equal(result[:], out1[:]))
}

func TestGateBuilder(t *testing.T) {
	// 1 \
	//     AND
	// 2 /      \
	//            5
	//              XOR -> 7
	//            6
	// 3 \      /
	//     OR
	// 4 /

	and := &AbstractGate{Type: GateAnd}
	or := &AbstractGate{Type: GateOr}
	xor := &AbstractGate{X: and, Y: or, Type: GateXor}

	circuit := &Circuit{
		Inputs:  []*AbstractGate{and, or},
		Outputs: []*AbstractGate{xor},
		results: nil,
	}

	out0 := Label{}

	out1 := Label{}
	out1[0] = 1

	circuit.Build([]Label{out0, out1})
	res := circuit.Run([]Label{and.Instance.(*GateTable).X1.L, and.Instance.(*GateTable).Y1.L, or.Instance.(*GateTable).X0.L, or.Instance.(*GateTable).Y1.L})
	fmt.Println(bytes.Equal(res[0][:], out0[:]))
}
