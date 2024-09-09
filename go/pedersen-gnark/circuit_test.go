// Package pedersen_gnark
// Copyright 2024 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen_gnark

import (
	"crypto/rand"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"math/big"
	"testing"
)

var curve = bn254.GetEdwardsCurve()

func getPoint() *bn254.PointAffine {
	k, err := rand.Int(rand.Reader, &curve.Order)
	if err != nil {
		panic(err)
	}

	defer func() { k = nil }()

	return new(bn254.PointAffine).ScalarMultiplication(&curve.Base, k)
}

func TestCircuitProveAndVerify(t *testing.T) {
	var circuit Circuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	pk, vk, err := groth16.Setup(r1cs)

	var assignment Circuit

	H := getPoint()
	G := getPoint()

	assignment.H = twistededwards.Point{
		X: H.X,
		Y: H.Y,
	}

	assignment.G = twistededwards.Point{
		X: G.X,
		Y: G.Y,
	}

	amount := big.NewInt(10)
	randomness := big.NewInt(20)

	assignment.Amount = new(fr.Element).SetBigInt(amount)
	assignment.Randomness = new(fr.Element).SetBigInt(randomness)

	commitment := new(bn254.PointAffine).Add(
		new(bn254.PointAffine).ScalarMultiplication(H, amount),
		new(bn254.PointAffine).ScalarMultiplication(G, randomness),
	)

	assignment.Commitment = twistededwards.Point{
		X: commitment.X,
		Y: commitment.Y,
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		panic(err)
	}

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		panic(err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
}
