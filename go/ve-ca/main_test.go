package ve_ca

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func TestEverythingWorks(t *testing.T) {
	g := MustRandomF()
	GenG := MustRandomG()
	GenH := MustRandomG()

	receiverPrv := MustRandomF()
	receiverShare := FPow(g, receiverPrv)

	v := MustRandomF()
	r := MustRandomF()

	proof, C := Encrypt(v, r, receiverShare, g, GenG, GenH)

	if !Verify(proof, receiverShare, C, GenG, GenH) {
		panic("failed to verify")
	}

	v_, r_ := Decrypt(receiverPrv, proof, C, GenG, GenH)

	if !bytes.Equal(FBytes(v), FBytes(v_)) {
		panic("invalid v recover")
	}

	if !bytes.Equal(FBytes(r), FBytes(r_)) {
		panic("invalid r recover")
	}
}

func TestEverythingWorksTimes(t *testing.T) {
	for testCase := range 1000 {
		t.Run(fmt.Sprintf("Instance simple #%d", testCase), TestEverythingWorks)
		t.Run(fmt.Sprintf("Instance encryptions #%d", testCase), TestCorruptedEncryptions)
		t.Run(fmt.Sprintf("Instance commitment #%d", testCase), TestCorruptedCommitment)
		t.Run(fmt.Sprintf("Instance malicious #%d", testCase), TestMaliciousEncryption)
	}
}

func TestCorruptedEncryptions(t *testing.T) {
	g := MustRandomF()
	GenG := MustRandomG()
	GenH := MustRandomG()

	receiverPrv := MustRandomF()
	receiverShare := FPow(g, receiverPrv)

	v := MustRandomF()
	r := MustRandomF()

	proof, C := Encrypt(v, r, receiverShare, g, GenG, GenH)

	proof.E[0][14] = MustRandomF()
	proof.E[1][14] = MustRandomF()

	if Verify(proof, receiverShare, C, GenG, GenH) {
		panic("verification should fail")
	}
}

func TestCorruptedCommitment(t *testing.T) {
	g := MustRandomF()
	GenG := MustRandomG()
	GenH := MustRandomG()

	receiverPrv := MustRandomF()
	receiverShare := FPow(g, receiverPrv)

	v := MustRandomF()
	r := MustRandomF()

	proof, C := Encrypt(v, r, receiverShare, g, GenG, GenH)

	C = MustRandomG()

	if Verify(proof, receiverShare, C, GenG, GenH) {
		panic("verification should fail")
	}
}

func TestMaliciousEncryption(t *testing.T) {
	g := MustRandomF()
	GenG := MustRandomG()
	GenH := MustRandomG()

	receiverPrv := MustRandomF()
	receiverShare := FPow(g, receiverPrv)

	v := MustRandomF()
	r := MustRandomF()

	proof, C := EncryptMalicious(v, r, receiverShare, g, GenG, GenH)

	C = MustRandomG()

	if Verify(proof, receiverShare, C, GenG, GenH) {
		panic("verification should fail")
	}
}

func EncryptMalicious(v, r F, s_r, g F, GenG, GenH G) (Proof, G) {
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

	rho0 := Hash(GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(X[:]), []byte{0x0})
	rho1 := Hash(GBytes(GenG), GBytes(GenH), GBytes(C), GArrBytes(X[:]), []byte{0x1})

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

		// Imitate malicious behaviour
		b := rand.Int() % 2
		e[b][i] = MustRandomF()
		e[b+2][i] = MustRandomF()
	}

	c := Hash(FBytes(rho0), FBytes(rho1), FArrBytes(e[0][:]), FArrBytes(e[1][:]), FArrBytes(e[2][:]), FArrBytes(e[3][:]))

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
