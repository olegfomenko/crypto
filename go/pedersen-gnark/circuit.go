// Package pedersen_gnark
// Copyright 2024 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen_gnark

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
)

type Circuit struct {
	Amount     frontend.Variable
	Randomness frontend.Variable
	Commitment twistededwards.Point `gnark:",public"`
	H          twistededwards.Point `gnark:",public"`
	G          twistededwards.Point `gnark:",public"`
}

func (c *Circuit) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	recoveredCommitment := curve.Add(curve.ScalarMul(c.H, c.Amount), curve.ScalarMul(c.G, c.Randomness))
	curve.API().AssertIsEqual(recoveredCommitment.X, c.Commitment.X)
	curve.API().AssertIsEqual(recoveredCommitment.Y, c.Commitment.Y)
	return nil
}
