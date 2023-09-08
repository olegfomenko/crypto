// Package pedersen
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"

	eth "github.com/ethereum/go-ethereum/crypto"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Hash function that should return the value in Curve.N field
var Hash func(...[]byte) *big.Int = defaultHash

// defaultHash - default hash function Keccak256
func defaultHash(bytes ...[]byte) *big.Int {
	data := make([][]byte, 0, len(bytes))
	for _, b := range bytes {
		data = append(data, uint256Bytes(b))
	}

	return new(big.Int).Mod(new(big.Int).SetBytes(eth.Keccak256(data...)), bn256.Order)
}

type Proof struct {
	E0 *big.Int
	C  []*bn256.G1
	S  []*big.Int
	N  int
}

// PedersenCommitment creates *bn256.G1 with pedersen commitment aH + rG
func PedersenCommitment(a, r *big.Int) *bn256.G1 {
	return Add(ScalarMul(H, a), ScalarMul(G, r))
}

// VerifyPedersenCommitment - verifies proof that C commitment commits the value in [0..2^n-1]
func VerifyPedersenCommitment(C *bn256.G1, proof Proof) error {
	var R []*bn256.G1

	for i := 0; i < proof.N; i++ {
		//calculating ei = Hash(si*G - e0(Ci - 2^i*H))
		siG := ScalarMul(G, proof.S[i])

		p := ScalarMul(H, pow2(i))
		p = Sub(proof.C[i], p)
		p = ScalarMul(p, proof.E0)
		p = Sub(siG, p)

		ei := hashPoints(p)

		R = append(R, ScalarMul(proof.C[i], ei))
	}

	// eo_ = Hash(Ro||R1||...Rn-1)
	e0_ := hashPoints(R...)

	// C = sum(Ci)
	Com := proof.C[0]
	for i := 1; i < proof.N; i++ {
		Com = Add(Com, proof.C[i])
	}

	if e0_.Cmp(proof.E0) != 0 {
		return errors.New("e0 != e0_")
	}
	if !bytes.Equal(C.Marshal(), Com.Marshal()) {
		return errors.New("C != sum(Ci)")
	}

	return nil
}

// CreatePedersenCommitment - creates Pedersen commitment for given val, and
// generates proof that given val lies in [0..2^n-1].
// Returns Proof, generated commitment and private key in case of success generation.
func CreatePedersenCommitment(val uint64, n int) (Proof, *bn256.G1, *big.Int, error) {
	// Converting into bit representation
	bitsStr := strconv.FormatUint(val, 2)
	var bits []bool
	for i := len(bitsStr) - 1; i >= 0; i-- {
		bits = append(bits, bitsStr[i] == '1')
	}

	if len(bits) > n {
		return Proof{}, nil, nil, errors.New("invalid value: greater then 2^n - 1")
	}

	// Adding leading zeros
	for len(bits) < n {
		bits = append(bits, false)
	}

	prv := big.NewInt(0)
	var r []*big.Int
	var k []*big.Int

	var R []*bn256.G1
	var C []*bn256.G1

	for i := 0; i < n; i++ {
		if bits[i] {
			ri, err := rand.Int(rand.Reader, bn256.Order)
			if err != nil {
				return Proof{}, nil, nil, err
			}
			prv = add(prv, ri)
			r = append(r, ri)

			// Ci = Com(2^i, ri)
			Ci := PedersenCommitment(pow2(i), ri)
			C = append(C, Ci)

			ki, err := rand.Int(rand.Reader, bn256.Order)
			if err != nil {
				return Proof{}, nil, nil, err
			}
			k = append(k, ki)

			// Hash(ki*G)
			kiG := ScalarMul(G, ki)
			ei := hashPoints(kiG)

			// Ri = Hash(ki*G)*Ci

			R = append(R, ScalarMul(Ci, ei))
			continue
		}

		ki0, err := rand.Int(rand.Reader, bn256.Order)
		if err != nil {
			return Proof{}, nil, nil, err
		}
		k = append(k, ki0)

		// Ri = ki0*G
		R = append(R, ScalarMul(G, ki0))

		// will be initialized later
		C = append(C, nil)
		// just placing nil value to be able to get corresponding r[i] for bit == 1 in future
		r = append(r, nil)
	}

	// eo = Hash(Ro||R1||...Rn-1)
	e0 := hashPoints(R...)

	var s []*big.Int

	for i := 0; i < n; i++ {
		if bits[i] {
			// si = ki + e0*r^i
			si := add(k[i], mul(e0, r[i]))
			s = append(s, si)
			continue
		}

		ki, err := rand.Int(rand.Reader, bn256.Order)
		if err != nil {
			return Proof{}, nil, nil, err
		}

		// ei = Hash(ki*G + e0*2^i*H)
		ei := hashPoints(PedersenCommitment(mul(e0, pow2(i)), ki))

		// Ci = Ri /ei = (ki0/ei)*G
		ei_inverse := new(big.Int).ModInverse(ei, bn256.Order)
		C[i] = ScalarMul(R[i], ei_inverse)

		prv = add(prv, mul(k[i], ei_inverse))

		// si = ki + (ki0 * e0)/ei
		si := add(ki, mul(mul(k[i], e0), ei_inverse))
		s = append(s, si)
	}

	Com := C[0]
	for i := 1; i < n; i++ {
		Com = Add(Com, C[i])
	}

	return Proof{
			E0: e0,
			C:  C,
			S:  s,
			N:  n,
		},
		Com,
		prv,
		nil
}

func add(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(x, y), bn256.Order)
}

func mul(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(x, y), bn256.Order)
}

func pow2(i int) *big.Int {
	return new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), bn256.Order)
}

func minus(val *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(val, big.NewInt(-1)), bn256.Order)
}

func hashPoints(points ...*bn256.G1) *big.Int {
	var data [][]byte
	for _, p := range points {
		data = append(data, X(p).Bytes())
		data = append(data, Y(p).Bytes())
	}

	return Hash(data...)
}

func uint256Bytes(val []byte) []byte {
	for len(val) < 32 {
		val = append([]byte{0}, val...)
	}
	return val
}
