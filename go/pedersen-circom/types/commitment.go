// Package types
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package types

import (
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type PrivateKey struct {
	R *BigInt `json:"R"`
	A *BigInt `json:"A"`
}

func NewPrivateKeyFromInput(input *InputData) *PrivateKey {
	return &PrivateKey{
		R: input.R,
		A: input.A,
	}
}

type Commitment struct {
	Public  *babyjub.Point `json:"public"`
	Private *PrivateKey    `json:"private"`
}

func NewCommitment(key *PrivateKey) *Commitment {
	// Com = aH + rG
	rG := babyjub.NewPoint().Mul(key.R.Int, G).Projective()
	aH := babyjub.NewPoint().Mul(key.A.Int, H).Projective()

	res := babyjub.NewPointProjective().Add(
		rG,
		aH,
	)

	return &Commitment{
		Public:  res.Affine(),
		Private: key,
	}
}

func NewCommitmentKeyFromInput(input *InputData) *Commitment {
	return NewCommitment(NewPrivateKeyFromInput(input))
}

func (c *Commitment) PublicKey() *babyjub.Point {
	return babyjub.NewPoint().Mul(c.Private.R.Int, G)
}
