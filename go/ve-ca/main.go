package ve_ca

import (
	"bytes"
	"fmt"
)

const k = 128

type Proof struct {
	E     [4][k]F
	Alpha [2][k]F
	P     [2][k]F
	X     [k]G
	Share [2][k]F
}

func Encrypt(v, r F, s_r, g F, GenG, GenH G) (Proof, G) {
	var p [4][k]F
	var x [2][k]F

	for i := range p {
		copy(p[i][:], MustSampleRandomF(k)[:])
	}

	for i := range x {
		copy(x[i][:], MustSampleRandomF(k)[:])
	}

	var s [4][k]F
	for i := range s {
		for j := range s[i] {
			s[i][j] = FPow(s_r, p[i][j])
		}
	}

	var X [k]G
	for i := range X {
		X[i] = GAdd(GMul(x[0][i], GenH), GMul(x[1][i], GenG))
	}

	C := GAdd(GMul(v, GenH), GMul(r, GenG))

	rho0 := Hash(
		GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(X[:]), []byte{0x0},
	)
	rho1 := Hash(
		GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(X[:]), []byte{0x1},
	)

	var alpha [4][k]F
	var e [4][k]F

	for i := range k {
		alpha[0][i] = FSub(x[0][i], FMul(rho0, v))
		alpha[1][i] = FSub(x[0][i], FMul(rho1, v))
		alpha[2][i] = FSub(x[1][i], FMul(rho0, r))
		alpha[3][i] = FSub(x[1][i], FMul(rho1, r))

		e[0][i] = E(s[0][i], alpha[0][i])
		e[1][i] = E(s[1][i], alpha[1][i])
		e[2][i] = E(s[2][i], alpha[2][i])
		e[3][i] = E(s[3][i], alpha[3][i])
	}

	c := Hash(
		FBytes(rho0),
		FBytes(rho1),
		FArrBytes(e[0][:]),
		FArrBytes(e[1][:]),
		FArrBytes(e[2][:]),
		FArrBytes(e[3][:]),
	)

	proof := Proof{
		E: e,
		X: X,
	}

	for i := range k {
		b := FBit(c, i)
		proof.Alpha[0][i] = alpha[b][i]
		proof.Alpha[1][i] = alpha[b+2][i]
		proof.P[0][i] = p[b][i]
		proof.P[1][i] = p[b+2][i]
		d := 1 - b
		proof.Share[0][i] = FPow(g, p[d][i])
		proof.Share[1][i] = FPow(g, p[d+2][i])
	}

	return proof, C
}

func Verify(proof Proof, s_r F, C, GenG, GenH G) bool {
	rho := []F{
		Hash(GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(proof.X[:]), []byte{0x0}),
		Hash(GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(proof.X[:]), []byte{0x1}),
	}

	c := Hash(FBytes(rho[0]), FBytes(rho[1]), FArrBytes(proof.E[0][:]), FArrBytes(proof.E[1][:]), FArrBytes(proof.E[2][:]), FArrBytes(proof.E[3][:]))

	for i := range k {
		b := FBit(c, i)
		X := GAdd(
			GAdd(
				GMul(proof.Alpha[0][i], GenH),
				GMul(proof.Alpha[1][i], GenG),
			), GMul(rho[b], C),
		)
		if !bytes.Equal(GBytes(X), GBytes(proof.X[i])) {
			fmt.Println(i)
			fmt.Println("invalid sigma answer")
			return false
		}

		e0 := E(FPow(s_r, proof.P[0][i]), proof.Alpha[0][i])
		e1 := E(FPow(s_r, proof.P[1][i]), proof.Alpha[1][i])

		if !bytes.Equal(FBytes(e0), FBytes(proof.E[b][i])) {
			fmt.Println(i)
			fmt.Println("invalid e0 answer")
			return false
		}

		if !bytes.Equal(FBytes(e1), FBytes(proof.E[b+2][i])) {
			fmt.Println(i)
			fmt.Println("invalid e1 answer")
			return false
		}
	}

	return true
}

func Decrypt(u F, proof Proof, C, GenG, GenH G) (F, F) {
	rho := [2]F{
		Hash(
			GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(proof.X[:]), []byte{0x0},
		),
		Hash(
			GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(proof.X[:]), []byte{0x1},
		),
	}

	c := Hash(
		FBytes(rho[0]),
		FBytes(rho[1]),
		FArrBytes(proof.E[0][:]),
		FArrBytes(proof.E[1][:]),
		FArrBytes(proof.E[2][:]),
		FArrBytes(proof.E[3][:]),
	)

	for i := range k {
		b := FBit(c, i)
		d := 1 - b

		keyV := FPow(proof.Share[0][i], u)
		keyR := FPow(proof.Share[1][i], u)

		alpha0 := D(keyV, proof.E[d][i])
		alpha1 := D(keyR, proof.E[d+2][i])

		v := FDiv(FSub(alpha0, proof.Alpha[0][i]), FSub(rho[b], rho[d]))
		r := FDiv(FSub(alpha1, proof.Alpha[1][i]), FSub(rho[b], rho[d]))

		C_ := GAdd(GMul(v, GenH), GMul(r, GenG))

		if bytes.Equal(GBytes(C_), GBytes(C)) {
			return v, r
		}
	}

	panic("Failed to recover")
}
